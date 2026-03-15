package test

import (
	"os"
	"strconv"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/freddierice/go-losetup/v2"
	"github.com/stretchr/testify/require"
)

func TestMirrorTarget(t *testing.T) {
	name := "test.mirrortarget"
	uuid := "e1f2a3b4-1234-5678-9abc-def012345678"

	dir := t.TempDir()
	numDevices := 2
	size := uint64(1024) * devmapper.SectorSize

	loops := make([]losetup.Device, numDevices)
	devices := make([]devmapper.MirrorDevice, numDevices)
	for i := range numDevices {
		backingFile := dir + "/mirror." + strconv.Itoa(i)
		f, err := os.Create(backingFile)
		require.NoError(t, err)
		require.NoError(t, f.Truncate(int64(size)))
		f.Close()

		loop, err := losetup.Attach(backingFile, 0, false)
		require.NoError(t, err)
		defer loop.Detach()
		loops[i] = loop

		devices[i] = devmapper.MirrorDevice{
			Device: loop.Path(),
			Offset: 0,
		}
	}

	m := devmapper.MirrorTable{
		Length:  size,
		LogType: "core",
		LogArgs: []string{"128"},
		Devices: devices,
	}
	require.NoError(t, devmapper.CreateAndLoad(name, uuid, 0, m))
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

func TestUserspaceMirrorTarget(t *testing.T) {
	dir := t.TempDir()
	numDevices := 2
	size := uint64(10) * devmapper.SectorSize

	files := make([]string, numDevices)
	devices := make([]devmapper.MirrorDevice, numDevices)
	for i := range numDevices {
		backingFile := dir + "/mirror." + strconv.Itoa(i)
		f, err := os.Create(backingFile)
		require.NoError(t, err)
		require.NoError(t, f.Truncate(int64(size)))
		f.Close()

		files[i] = backingFile
		devices[i] = devmapper.MirrorDevice{
			Device: backingFile,
			Offset: 0,
		}
	}

	m := devmapper.MirrorTable{
		Length:  size,
		LogType: "core",
		LogArgs: []string{"128"},
		Devices: devices,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, m)
	require.NoError(t, err)
	defer v.Close()

	// Write data
	writeBuf := make([]byte, devmapper.SectorSize)
	copy(writeBuf, "Hello, mirror!!!")
	n, err := v.WriteAt(writeBuf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)

	// Read back
	readBuf := make([]byte, devmapper.SectorSize)
	n, err = v.ReadAt(readBuf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)
	require.Equal(t, writeBuf, readBuf)

	v.Close()

	// Verify both backing files have the same data (mirrored)
	for i := range numDevices {
		data, err := os.ReadFile(files[i])
		require.NoError(t, err)
		expected := make([]byte, size)
		copy(expected, "Hello, mirror!!!")
		require.Equal(t, expected, data, "mirror device %d data mismatch", i)
	}
}

func TestUserspaceMirrorTargetWithOffset(t *testing.T) {
	dir := t.TempDir()
	size := uint64(10) * devmapper.SectorSize
	offset := uint64(5) * devmapper.SectorSize

	files := make([]string, 2)
	devices := make([]devmapper.MirrorDevice, 2)
	for i := range 2 {
		backingFile := dir + "/mirror." + strconv.Itoa(i)
		f, err := os.Create(backingFile)
		require.NoError(t, err)
		require.NoError(t, f.Truncate(int64(size+offset)))
		f.Close()

		files[i] = backingFile
		devices[i] = devmapper.MirrorDevice{
			Device: backingFile,
			Offset: offset,
		}
	}

	m := devmapper.MirrorTable{
		Length:  size,
		LogType: "core",
		LogArgs: []string{"128"},
		Devices: devices,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, m)
	require.NoError(t, err)

	writeBuf := make([]byte, devmapper.SectorSize)
	copy(writeBuf, "offset mirror!!!")
	_, err = v.WriteAt(writeBuf, 0)
	require.NoError(t, err)
	v.Close()

	// Verify data at the correct offset in backing files
	for i := range 2 {
		f, err := os.Open(files[i])
		require.NoError(t, err)
		buf := make([]byte, devmapper.SectorSize)
		_, err = f.ReadAt(buf, int64(offset))
		require.NoError(t, err)
		require.Equal(t, writeBuf, buf, "mirror device %d offset data mismatch", i)
		f.Close()
	}
}
