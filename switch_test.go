package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSwitchBuildSpec(t *testing.T) {
	t.Parallel()

	s := SwitchTable{
		Start:      0,
		Length:     10000 * SectorSize,
		RegionSize: 128 * SectorSize,
		Devices: []SwitchDevice{
			{Device: "/dev/sda1", Offset: 0},
			{Device: "/dev/sdb1", Offset: 0},
		},
	}
	require.Equal(t, "2 128 /dev/sda1 0 /dev/sdb1 0", s.buildSpec())
	require.Equal(t, "switch", s.targetType())
	require.Equal(t, uint64(0), s.start())
	require.Equal(t, uint64(10000*SectorSize), s.length())
}
