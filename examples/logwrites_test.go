package test

import (
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/freddierice/go-losetup/v2"
	"github.com/stretchr/testify/require"
)

func TestLogWritesTarget(t *testing.T) {
	name := "test.logwritestarget"
	uuid := "f1e2d3c4-9999-8888-7777-666655554444"

	dir := t.TempDir()
	size := uint64(10) * devmapper.SectorSize

	backingFile := dir + "/data"
	f, err := os.Create(backingFile)
	require.NoError(t, err)
	text := "Hello, logwrites"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	logFile := dir + "/log"
	f, err = os.Create(logFile)
	require.NoError(t, err)
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	dataLoop, err := losetup.Attach(backingFile, 0, false)
	require.NoError(t, err)
	defer dataLoop.Detach()

	logLoop, err := losetup.Attach(logFile, 0, false)
	require.NoError(t, err)
	defer logLoop.Detach()

	l := devmapper.LogWritesTable{
		Length:    size,
		Device:    dataLoop.Path(),
		LogDevice: logLoop.Path(),
	}
	require.NoError(t, devmapper.CreateAndLoad(name, uuid, 0, l))
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

func TestUserspaceLogWritesTarget(t *testing.T) {
	dir := t.TempDir()
	size := uint64(10) * devmapper.SectorSize

	backingFile := dir + "/data"
	f, err := os.Create(backingFile)
	require.NoError(t, err)
	text := "Hello, logwrites"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	logFile := dir + "/log"
	f, err = os.Create(logFile)
	require.NoError(t, err)
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	l := devmapper.LogWritesTable{
		Length:    size,
		Device:    backingFile,
		LogDevice: logFile,
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, l)
	require.NoError(t, err)
	defer v.Close()

	// Read
	buf := make([]byte, devmapper.SectorSize)
	n, err := v.ReadAt(buf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)
	expected := make([]byte, devmapper.SectorSize)
	copy(expected, text)
	require.Equal(t, expected, buf)

	// Write
	writeBuf := make([]byte, devmapper.SectorSize)
	copy(writeBuf, "logwrites write!")
	n, err = v.WriteAt(writeBuf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)

	// Read back
	readBuf := make([]byte, devmapper.SectorSize)
	_, err = v.ReadAt(readBuf, 0)
	require.NoError(t, err)
	require.Equal(t, writeBuf, readBuf)
}
