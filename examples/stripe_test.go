package test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/freddierice/go-losetup/v2"
	"github.com/stretchr/testify/require"
)

func TestStripeTarget(t *testing.T) {
	name := "test.stripetarget"
	uuid := "d1e2f3a4-5555-6666-7777-888899990001"

	dir := t.TempDir()
	numDevices := 3
	chunkSize := uint64(2) * devmapper.SectorSize // 1024 bytes
	sectorsPerDevice := uint64(6)

	loops := make([]losetup.Device, numDevices)
	devices := make([]devmapper.StripeDevice, numDevices)
	for i := range numDevices {
		backingFile := dir + "/stripe." + strconv.Itoa(i)
		f, err := os.Create(backingFile)
		require.NoError(t, err)
		require.NoError(t, f.Truncate(int64(sectorsPerDevice*devmapper.SectorSize)))
		f.Close()

		loop, err := losetup.Attach(backingFile, 0, false)
		require.NoError(t, err)
		defer loop.Detach()
		loops[i] = loop

		devices[i] = devmapper.StripeDevice{
			Device: loop.Path(),
			Offset: 0,
		}
	}

	// Total size = numDevices * sectorsPerDevice * SectorSize
	totalSize := uint64(numDevices) * sectorsPerDevice * devmapper.SectorSize
	s := devmapper.StripeTable{
		Length:    totalSize,
		ChunkSize: chunkSize,
		Devices:  devices,
	}
	require.NoError(t, devmapper.CreateAndLoad(name, uuid, 0, s))
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

	mapper := "/dev/mapper/" + name
	require.NoError(t, waitForFile(mapper))

	// Write data via the mapper device
	f, err := os.OpenFile(mapper, os.O_WRONLY, 0)
	require.NoError(t, err)
	data := make([]byte, totalSize)
	for i := range data {
		data[i] = byte(i % 251) // use prime to avoid alignment coincidences
	}
	_, err = f.Write(data)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	// Read it back
	readData, err := os.ReadFile(mapper)
	require.NoError(t, err)
	require.Equal(t, data, readData)
}

func TestUserspaceStripeReadWrite(t *testing.T) {
	dir := t.TempDir()
	numDevices := 3
	chunkSize := uint64(2) * devmapper.SectorSize
	sectorsPerDevice := uint64(6)

	devices := make([]devmapper.StripeDevice, numDevices)
	for i := range numDevices {
		backingFile := dir + "/stripe." + strconv.Itoa(i)
		f, err := os.Create(backingFile)
		require.NoError(t, err)
		require.NoError(t, f.Truncate(int64(sectorsPerDevice*devmapper.SectorSize)))
		f.Close()

		devices[i] = devmapper.StripeDevice{
			Device: backingFile,
			Offset: 0,
		}
	}

	totalSize := uint64(numDevices) * sectorsPerDevice * devmapper.SectorSize
	s := devmapper.StripeTable{
		Length:    totalSize,
		ChunkSize: chunkSize,
		Devices:  devices,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, s)
	require.NoError(t, err)
	defer v.Close()

	// Write data
	data := make([]byte, totalSize)
	for i := range data {
		data[i] = byte(i % 251)
	}
	n, err := v.WriteAt(data, 0)
	require.NoError(t, err)
	require.Equal(t, int(totalSize), n)

	// Read it back
	buf := make([]byte, totalSize)
	n, err = v.ReadAt(buf, 0)
	require.NoError(t, err)
	require.Equal(t, int(totalSize), n)
	require.Equal(t, data, buf)

	// Verify stripe layout on disk: first chunk goes to device 0,
	// second chunk to device 1, third to device 2, etc.
	for i := range numDevices {
		backingFile := dir + "/stripe." + strconv.Itoa(i)
		raw, err := os.ReadFile(backingFile)
		require.NoError(t, err)

		// Each device should have sectorsPerDevice sectors.
		// For device i, it gets chunks i, i+numDevices, i+2*numDevices, ...
		for chunkOnDev := uint64(0); chunkOnDev < sectorsPerDevice*devmapper.SectorSize/chunkSize; chunkOnDev++ {
			virtualChunk := chunkOnDev*uint64(numDevices) + uint64(i)
			virtualOff := virtualChunk * chunkSize
			diskOff := chunkOnDev * chunkSize
			expected := data[virtualOff : virtualOff+chunkSize]
			actual := raw[diskOff : diskOff+chunkSize]
			require.Equal(t, expected, actual, fmt.Sprintf("stripe layout mismatch for device %d chunk %d", i, chunkOnDev))
		}
	}
}

func TestUserspaceStripePartialRead(t *testing.T) {
	dir := t.TempDir()
	numDevices := 2
	chunkSize := uint64(2) * devmapper.SectorSize
	sectorsPerDevice := uint64(4)

	devices := make([]devmapper.StripeDevice, numDevices)
	for i := range numDevices {
		backingFile := dir + "/stripe." + strconv.Itoa(i)
		f, err := os.Create(backingFile)
		require.NoError(t, err)
		require.NoError(t, f.Truncate(int64(sectorsPerDevice*devmapper.SectorSize)))
		f.Close()

		devices[i] = devmapper.StripeDevice{
			Device: backingFile,
			Offset: 0,
		}
	}

	totalSize := uint64(numDevices) * sectorsPerDevice * devmapper.SectorSize
	s := devmapper.StripeTable{
		Length:    totalSize,
		ChunkSize: chunkSize,
		Devices:  devices,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, s)
	require.NoError(t, err)
	defer v.Close()

	// Write data
	data := make([]byte, totalSize)
	for i := range data {
		data[i] = byte(i % 199)
	}
	_, err = v.WriteAt(data, 0)
	require.NoError(t, err)

	// Read from the middle, crossing a chunk boundary
	readOff := int64(chunkSize - devmapper.SectorSize)
	readLen := 2 * devmapper.SectorSize
	buf := make([]byte, readLen)
	n, err := v.ReadAt(buf, readOff)
	require.NoError(t, err)
	require.Equal(t, readLen, n)
	require.Equal(t, data[readOff:readOff+int64(readLen)], buf)
}
