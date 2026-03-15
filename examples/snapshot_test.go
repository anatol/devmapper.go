package test

import (
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/freddierice/go-losetup/v2"
	"github.com/stretchr/testify/require"
)

func TestUserspaceSnapshotOriginTarget(t *testing.T) {
	dir := t.TempDir()
	backingFile := dir + "/origin"
	f, err := os.Create(backingFile)
	require.NoError(t, err)

	text := "Hello, origin!!!"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	size := uint64(10) * devmapper.SectorSize
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	s := devmapper.SnapshotOriginTable{
		Length:       size,
		OriginDevice: backingFile,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, s)
	require.NoError(t, err)
	defer v.Close()

	buf := make([]byte, devmapper.SectorSize)
	n, err := v.ReadAt(buf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)

	expected := make([]byte, devmapper.SectorSize)
	copy(expected, text)
	require.Equal(t, expected, buf)

	// Test write
	writeBuf := make([]byte, devmapper.SectorSize)
	copy(writeBuf, "origin write!!!!")
	n, err = v.WriteAt(writeBuf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)

	readBuf := make([]byte, devmapper.SectorSize)
	_, err = v.ReadAt(readBuf, 0)
	require.NoError(t, err)
	require.Equal(t, writeBuf, readBuf)
}

func TestSnapshotOriginTarget(t *testing.T) {
	dir := t.TempDir()
	backingFile := dir + "/origin"
	f, err := os.Create(backingFile)
	require.NoError(t, err)
	text := "Hello, origin!!!"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	size := uint64(10) * devmapper.SectorSize
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	loop, err := losetup.Attach(backingFile, 0, false)
	require.NoError(t, err)
	defer loop.Detach()

	name := "test.snapshotorigin"
	uuid := "a1b2c3d4-aaaa-bbbb-cccc-111122223333"
	s := devmapper.SnapshotOriginTable{
		Length:       size,
		OriginDevice: loop.Path(),
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

	data, err := os.ReadFile(mapper)
	require.NoError(t, err)
	expected := make([]byte, size)
	copy(expected, text)
	require.Equal(t, expected, data)
}
