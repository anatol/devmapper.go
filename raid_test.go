package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRaidBuildSpec(t *testing.T) {
	t.Parallel()

	r := RaidTable{
		Start:    0,
		Length:   10000 * SectorSize,
		RaidType: "raid1",
		Params:   []string{"3", "rebuild", "/dev/sdb1", "128"},
		Devices: []RaidDevice{
			{MetaDevice: "-", DataDevice: "/dev/sda1"},
			{MetaDevice: "-", DataDevice: "/dev/sdb1"},
		},
	}
	require.Equal(t, "raid1 4 3 rebuild /dev/sdb1 128 2 - /dev/sda1 - /dev/sdb1", r.buildSpec())
	require.Equal(t, "raid", r.targetType())
	require.Equal(t, uint64(0), r.start())
	require.Equal(t, uint64(10000*SectorSize), r.length())
}

func TestRaidBuildSpecRaid0(t *testing.T) {
	t.Parallel()

	r := RaidTable{
		Start:    0,
		Length:   10000 * SectorSize,
		RaidType: "raid0",
		Params:   []string{"128"},
		Devices: []RaidDevice{
			{MetaDevice: "-", DataDevice: "/dev/sda1"},
			{MetaDevice: "-", DataDevice: "/dev/sdb1"},
			{MetaDevice: "-", DataDevice: "/dev/sdc1"},
		},
	}
	require.Equal(t, "raid0 1 128 3 - /dev/sda1 - /dev/sdb1 - /dev/sdc1", r.buildSpec())
}
