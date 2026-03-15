package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlakeyBuildSpec(t *testing.T) {
	t.Parallel()

	f := FlakeyTable{
		Start:        0,
		Length:       1000 * SectorSize,
		Device:       "/dev/sda1",
		Offset:       0,
		UpInterval:   30,
		DownInterval: 5,
	}
	require.Equal(t, "/dev/sda1 0 30 5", f.buildSpec())
	require.Equal(t, "flakey", f.targetType())
	require.Equal(t, uint64(0), f.start())
	require.Equal(t, uint64(1000*SectorSize), f.length())
}

func TestFlakeyBuildSpecWithFeatures(t *testing.T) {
	t.Parallel()

	f := FlakeyTable{
		Start:        0,
		Length:       1000 * SectorSize,
		Device:       "/dev/sda1",
		Offset:       512 * SectorSize,
		UpInterval:   60,
		DownInterval: 10,
		Features:     []string{"drop_writes", "error_writes"},
	}
	require.Equal(t, "/dev/sda1 512 60 10 2 drop_writes error_writes", f.buildSpec())
}
