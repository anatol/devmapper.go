package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDustBuildSpec(t *testing.T) {
	t.Parallel()

	d := DustTable{
		Start:     0,
		Length:    10000 * SectorSize,
		Device:    "/dev/sda1",
		Offset:    0,
		BlockSize: 8 * SectorSize,
	}
	require.Equal(t, "/dev/sda1 0 4096", d.buildSpec())
	require.Equal(t, "dust", d.targetType())
	require.Equal(t, uint64(0), d.start())
	require.Equal(t, uint64(10000*SectorSize), d.length())
}

func TestDustBuildSpecWithOffset(t *testing.T) {
	t.Parallel()

	d := DustTable{
		Start:     0,
		Length:    10000 * SectorSize,
		Device:    "/dev/sda1",
		Offset:    1024 * SectorSize,
		BlockSize: 8 * SectorSize,
	}
	require.Equal(t, "/dev/sda1 1024 4096", d.buildSpec())
}
