package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStripeBuildSpec(t *testing.T) {
	t.Parallel()

	s := StripeTable{
		Start:     0,
		Length:    1000 * SectorSize,
		ChunkSize: 128 * SectorSize,
		Devices: []StripeDevice{
			{Device: "/dev/sda1", Offset: 0},
			{Device: "/dev/sdb1", Offset: 0},
		},
	}
	require.Equal(t, "2 128 /dev/sda1 0 /dev/sdb1 0", s.buildSpec())
	require.Equal(t, "striped", s.targetType())
	require.Equal(t, uint64(0), s.start())
	require.Equal(t, uint64(1000*SectorSize), s.length())
}

func TestStripeBuildSpecWithOffset(t *testing.T) {
	t.Parallel()

	s := StripeTable{
		Start:     0,
		Length:    2000 * SectorSize,
		ChunkSize: 64 * SectorSize,
		Devices: []StripeDevice{
			{Device: "/dev/sda1", Offset: 512 * SectorSize},
			{Device: "/dev/sdb1", Offset: 1024 * SectorSize},
			{Device: "/dev/sdc1", Offset: 0},
		},
	}
	require.Equal(t, "3 64 /dev/sda1 512 /dev/sdb1 1024 /dev/sdc1 0", s.buildSpec())
}
