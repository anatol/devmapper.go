package test

import (
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/freddierice/go-losetup/v2"
	"github.com/stretchr/testify/require"
)

func TestDelayTarget(t *testing.T) {
	name := "test.delaytarget"
	uuid := "c1d2e3f4-aaaa-bbbb-cccc-ddddeeeeffff"

	dir := t.TempDir()
	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	require.NoError(t, err)

	text := "Hello, delay!!!"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	size := uint64(10) * devmapper.SectorSize
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	loop, err := losetup.Attach(backingFile, 0, false)
	require.NoError(t, err)
	defer loop.Detach()

	d := devmapper.DelayTable{
		Length:     size,
		ReadDevice: loop.Path(),
		ReadDelay:  0,
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

func TestUserspaceDelayTarget(t *testing.T) {
	dir := t.TempDir()
	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	require.NoError(t, err)

	text := "Hello, delay!!!"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	size := uint64(10) * devmapper.SectorSize
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	d := devmapper.DelayTable{
		Length:     size,
		ReadDevice: backingFile,
		ReadDelay:  100,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDONLY, 0, d)
	require.NoError(t, err)
	defer v.Close()

	buf := make([]byte, size)
	n, err := v.ReadAt(buf, 0)
	require.NoError(t, err)
	require.Equal(t, int(size), n)

	expected := make([]byte, size)
	copy(expected, text)
	require.Equal(t, expected, buf)
}

func TestUserspaceDelayTargetSeparateDevices(t *testing.T) {
	dir := t.TempDir()
	size := uint64(10) * devmapper.SectorSize

	readFile := dir + "/read"
	f, err := os.Create(readFile)
	require.NoError(t, err)
	_, err = f.WriteString("read data here!")
	require.NoError(t, err)
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	writeFile := dir + "/write"
	f, err = os.Create(writeFile)
	require.NoError(t, err)
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	d := devmapper.DelayTable{
		Length:      size,
		ReadDevice:  readFile,
		ReadDelay:   0,
		WriteDevice: writeFile,
		WriteDelay:  0,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, d)
	require.NoError(t, err)
	defer v.Close()

	// Read from readFile
	buf := make([]byte, devmapper.SectorSize)
	n, err := v.ReadAt(buf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)
	expected := make([]byte, devmapper.SectorSize)
	copy(expected, "read data here!")
	require.Equal(t, expected, buf)

	// Write goes to writeFile
	writeBuf := make([]byte, devmapper.SectorSize)
	copy(writeBuf, "written data!!!")
	n, err = v.WriteAt(writeBuf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)

	// Verify writeFile got the data
	raw, err := os.ReadFile(writeFile)
	require.NoError(t, err)
	expectedWrite := make([]byte, size)
	copy(expectedWrite, "written data!!!")
	require.Equal(t, expectedWrite, raw)
}
