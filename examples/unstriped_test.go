package test

import (
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/freddierice/go-losetup/v2"
	"github.com/stretchr/testify/require"
)

func TestUnstripedTarget(t *testing.T) {
	name := "test.unstripedtarget"
	uuid := "d1e2f3a4-5555-6666-7777-888899990010"

	dir := t.TempDir()
	numStripes := uint64(3)
	chunkSize := uint64(2) * devmapper.SectorSize
	chunksPerStripe := uint64(4)
	// The underlying device holds numStripes * chunksPerStripe chunks
	totalBackendSize := numStripes * chunksPerStripe * chunkSize

	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	require.NoError(t, err)
	require.NoError(t, f.Truncate(int64(totalBackendSize)))
	f.Close()

	loop, err := losetup.Attach(backingFile, 0, false)
	require.NoError(t, err)
	defer loop.Detach()

	stripeNum := uint64(1)
	unstripedSize := chunksPerStripe * chunkSize

	u := devmapper.UnstripedTable{
		Length:     unstripedSize,
		NumStripes: numStripes,
		ChunkSize:  chunkSize,
		StripeNum:  stripeNum,
		Device:     loop.Path(),
		Offset:     0,
	}
	require.NoError(t, devmapper.CreateAndLoad(name, uuid, 0, u))
	defer devmapper.Remove(name)

	got, err := devInfo(name)
	require.NoError(t, err)
	checkDevInfo(t, got, map[string]string{
		PropName:          name,
		PropTargetsNum:    "1",
		PropState:         "ACTIVE",
		PropTablesPresent: "LIVE",
		PropUUID:          uuid,
	})
}

func TestUserspaceUnstriped(t *testing.T) {
	dir := t.TempDir()
	numStripes := uint64(3)
	chunkSize := uint64(2) * devmapper.SectorSize
	chunksPerStripe := uint64(4)
	totalBackendSize := numStripes * chunksPerStripe * chunkSize

	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	require.NoError(t, err)

	// Fill the backing file with known data
	data := make([]byte, totalBackendSize)
	for i := range data {
		data[i] = byte(i % 251)
	}
	_, err = f.Write(data)
	require.NoError(t, err)
	f.Close()

	stripeNum := uint64(1)
	unstripedSize := chunksPerStripe * chunkSize

	u := devmapper.UnstripedTable{
		Length:     unstripedSize,
		NumStripes: numStripes,
		ChunkSize:  chunkSize,
		StripeNum:  stripeNum,
		Device:     backingFile,
		Offset:     0,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDONLY, 0, u)
	require.NoError(t, err)
	defer v.Close()

	buf := make([]byte, unstripedSize)
	n, err := v.ReadAt(buf, 0)
	require.NoError(t, err)
	require.Equal(t, int(unstripedSize), n)

	// Verify: chunk i of the unstriped view should correspond to
	// chunk (i*numStripes + stripeNum) of the backing file
	for chunkIdx := uint64(0); chunkIdx < chunksPerStripe; chunkIdx++ {
		virtualOff := chunkIdx * chunkSize
		backendOff := chunkIdx*numStripes*chunkSize + stripeNum*chunkSize
		expected := data[backendOff : backendOff+chunkSize]
		actual := buf[virtualOff : virtualOff+chunkSize]
		require.Equal(t, expected, actual, "unstriped chunk %d mismatch", chunkIdx)
	}
}

func TestUserspaceUnstripedWrite(t *testing.T) {
	dir := t.TempDir()
	numStripes := uint64(2)
	chunkSize := uint64(2) * devmapper.SectorSize
	chunksPerStripe := uint64(3)
	totalBackendSize := numStripes * chunksPerStripe * chunkSize

	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	require.NoError(t, err)
	require.NoError(t, f.Truncate(int64(totalBackendSize)))
	f.Close()

	stripeNum := uint64(0)
	unstripedSize := chunksPerStripe * chunkSize

	u := devmapper.UnstripedTable{
		Length:     unstripedSize,
		NumStripes: numStripes,
		ChunkSize:  chunkSize,
		StripeNum:  stripeNum,
		Device:     backingFile,
		Offset:     0,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, u)
	require.NoError(t, err)

	// Write pattern
	writeData := make([]byte, unstripedSize)
	for i := range writeData {
		writeData[i] = byte(i%199 + 1)
	}
	n, err := v.WriteAt(writeData, 0)
	require.NoError(t, err)
	require.Equal(t, int(unstripedSize), n)

	// Read back and verify
	buf := make([]byte, unstripedSize)
	n, err = v.ReadAt(buf, 0)
	require.NoError(t, err)
	require.Equal(t, int(unstripedSize), n)
	require.Equal(t, writeData, buf)

	v.Close()

	// Verify on-disk layout
	raw, err := os.ReadFile(backingFile)
	require.NoError(t, err)
	for chunkIdx := uint64(0); chunkIdx < chunksPerStripe; chunkIdx++ {
		virtualOff := chunkIdx * chunkSize
		backendOff := chunkIdx*numStripes*chunkSize + stripeNum*chunkSize
		expected := writeData[virtualOff : virtualOff+chunkSize]
		actual := raw[backendOff : backendOff+chunkSize]
		require.Equal(t, expected, actual, "unstriped write chunk %d mismatch", chunkIdx)
	}
}
