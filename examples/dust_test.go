package test

import (
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/freddierice/go-losetup/v2"
	"github.com/stretchr/testify/require"
)

func TestDustTarget(t *testing.T) {
	name := "test.dusttarget"
	uuid := "b1c2d3e4-dddd-eeee-ffff-111122223333"

	dir := t.TempDir()
	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	require.NoError(t, err)

	text := "Hello, dust!!!!!"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	size := uint64(10) * devmapper.SectorSize
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	loop, err := losetup.Attach(backingFile, 0, false)
	require.NoError(t, err)
	defer loop.Detach()

	d := devmapper.DustTable{
		Length:    size,
		Device:    loop.Path(),
		BlockSize: devmapper.SectorSize,
	}
	require.NoError(t, devmapper.CreateAndLoad(name, uuid, 0, d))
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

	data, err := os.ReadFile(mapper)
	require.NoError(t, err)
	expected := make([]byte, size)
	copy(expected, text)
	require.Equal(t, expected, data)
}

func TestUserspaceDustTarget(t *testing.T) {
	dir := t.TempDir()
	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	require.NoError(t, err)

	text := "Hello, dust!!!!!"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	size := uint64(10) * devmapper.SectorSize
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	d := devmapper.DustTable{
		Length:    size,
		Device:    backingFile,
		BlockSize: devmapper.SectorSize,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, d)
	require.NoError(t, err)
	defer v.Close()

	buf := make([]byte, size)
	n, err := v.ReadAt(buf, 0)
	require.NoError(t, err)
	require.Equal(t, int(size), n)

	expected := make([]byte, size)
	copy(expected, text)
	require.Equal(t, expected, buf)

	// Test write
	writeBuf := make([]byte, devmapper.SectorSize)
	copy(writeBuf, "dust write!!!!!")
	n, err = v.WriteAt(writeBuf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)

	readBuf := make([]byte, devmapper.SectorSize)
	n, err = v.ReadAt(readBuf, 0)
	require.NoError(t, err)
	require.Equal(t, writeBuf, readBuf)
}
