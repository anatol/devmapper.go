package test

import (
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/stretchr/testify/require"
)

func TestUserspaceMultipathTarget(t *testing.T) {
	dir := t.TempDir()
	size := uint64(10) * devmapper.SectorSize

	backingFile := dir + "/path0"
	f, err := os.Create(backingFile)
	require.NoError(t, err)
	text := "Hello, multipath"
	_, err = f.WriteString(text)
	require.NoError(t, err)
	require.NoError(t, f.Truncate(int64(size)))
	f.Close()

	m := devmapper.MultipathTable{
		Length:    size,
		HWHandler: "0",
		PathGroups: []devmapper.MultipathGroup{
			{
				PathSelector: "round-robin",
				SelectorArgs: []string{"0"},
				Paths: []devmapper.MultipathPath{
					{Device: backingFile},
				},
			},
		},
	}

	v, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, m)
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
	copy(writeBuf, "multipath write!")
	n, err = v.WriteAt(writeBuf, 0)
	require.NoError(t, err)
	require.Equal(t, devmapper.SectorSize, n)

	readBuf := make([]byte, devmapper.SectorSize)
	_, err = v.ReadAt(readBuf, 0)
	require.NoError(t, err)
	require.Equal(t, writeBuf, readBuf)
}
