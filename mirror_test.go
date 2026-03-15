package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMirrorBuildSpec(t *testing.T) {
	t.Parallel()

	m := MirrorTable{
		Start:   0,
		Length:  1000 * SectorSize,
		LogType: "core",
		LogArgs: []string{"1", "nosync"},
		Devices: []MirrorDevice{
			{Device: "/dev/sda1", Offset: 0},
			{Device: "/dev/sdb1", Offset: 0},
		},
		Features: []string{"handle_errors"},
	}
	require.Equal(t, "core 2 1 nosync 2 /dev/sda1 0 /dev/sdb1 0 1 handle_errors", m.buildSpec())
	require.Equal(t, "mirror", m.targetType())
	require.Equal(t, uint64(0), m.start())
	require.Equal(t, uint64(1000*SectorSize), m.length())
}

func TestMirrorBuildSpecNoFeatures(t *testing.T) {
	t.Parallel()

	m := MirrorTable{
		Start:   0,
		Length:  1000 * SectorSize,
		LogType: "disk",
		LogArgs: []string{"1", "/dev/sdc1"},
		Devices: []MirrorDevice{
			{Device: "/dev/sda1", Offset: 0},
			{Device: "/dev/sdb1", Offset: 0},
		},
	}
	require.Equal(t, "disk 2 1 /dev/sdc1 2 /dev/sda1 0 /dev/sdb1 0 0", m.buildSpec())
}
