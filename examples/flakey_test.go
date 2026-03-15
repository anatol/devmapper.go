package test

import (
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/freddierice/go-losetup/v2"
	"github.com/stretchr/testify/require"
)

func TestFlakeyTarget(t *testing.T) {
	name := "test.flakeytarget"
	uuid := "f1a2b3c4-1111-2222-3333-aabbccddeeff"

	dir := t.TempDir()
	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	require.NoError(t, err)

	text := "Hello, flakey!!"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	size := uint64(10) * devmapper.SectorSize
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	loop, err := losetup.Attach(backingFile, 0, false)
	require.NoError(t, err)
	defer loop.Detach()

	fl := devmapper.FlakeyTable{
		Length:       size,
		Device:       loop.Path(),
		UpInterval:   9999, // long up interval so test can complete
		DownInterval: 0,
	}
	require.NoError(t, devmapper.CreateAndLoad(name, uuid, 0, fl))
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

func TestUserspaceFlakeyTarget(t *testing.T) {
	dir := t.TempDir()
	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	require.NoError(t, err)

	text := "Hello, flakey!!"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	size := uint64(10) * devmapper.SectorSize
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	fl := devmapper.FlakeyTable{
		Length:       size,
		Device:       backingFile,
		UpInterval:   9999,
		DownInterval: 0,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, fl)
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
	copy(writeBuf, "flakey write!!!")
	n, err = v.WriteAt(writeBuf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)

	// Read back
	readBuf := make([]byte, devmapper.SectorSize)
	n, err = v.ReadAt(readBuf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)
	require.Equal(t, writeBuf, readBuf)
}
