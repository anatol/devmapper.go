package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultipathBuildSpec(t *testing.T) {
	t.Parallel()

	m := MultipathTable{
		Start:     0,
		Length:    10000 * SectorSize,
		Features:  []string{"queue_if_no_path"},
		HWHandler: "0",
		PathGroups: []MultipathGroup{
			{
				PathSelector: "round-robin",
				SelectorArgs: []string{"0"},
				Paths: []MultipathPath{
					{Device: "/dev/sda", Args: []string{"1"}},
					{Device: "/dev/sdb", Args: []string{"1"}},
				},
			},
		},
	}
	require.Equal(t, "1 queue_if_no_path 0 0 1 round-robin 1 0 2 /dev/sda 1 1 /dev/sdb 1 1", m.buildSpec())
	require.Equal(t, "multipath", m.targetType())
	require.Equal(t, uint64(0), m.start())
	require.Equal(t, uint64(10000*SectorSize), m.length())
}

func TestMultipathBuildSpecMultipleGroups(t *testing.T) {
	t.Parallel()

	m := MultipathTable{
		Start:     0,
		Length:    10000 * SectorSize,
		HWHandler: "0",
		PathGroups: []MultipathGroup{
			{
				PathSelector: "round-robin",
				SelectorArgs: []string{"0"},
				Paths: []MultipathPath{
					{Device: "/dev/sda"},
					{Device: "/dev/sdb"},
				},
			},
			{
				PathSelector: "round-robin",
				SelectorArgs: []string{"0"},
				Paths: []MultipathPath{
					{Device: "/dev/sdc"},
				},
			},
		},
	}
	require.Equal(t, "0 0 0 2 round-robin 1 0 2 /dev/sda 0 /dev/sdb 0 round-robin 1 0 1 /dev/sdc 0", m.buildSpec())
}
