package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestThinPoolBuildSpec(t *testing.T) {
	t.Parallel()

	tp := ThinPoolTable{
		Start:          0,
		Length:         10000 * SectorSize,
		MetadataDevice: "/dev/sda1",
		DataDevice:     "/dev/sdb1",
		DataBlockSize:  128 * SectorSize,
		LowWaterMark:  100,
	}
	require.Equal(t, "/dev/sda1 /dev/sdb1 128 100 0", tp.buildSpec())
	require.Equal(t, "thin-pool", tp.targetType())
	require.Equal(t, uint64(0), tp.start())
	require.Equal(t, uint64(10000*SectorSize), tp.length())
}

func TestThinPoolBuildSpecWithFeatures(t *testing.T) {
	t.Parallel()

	tp := ThinPoolTable{
		Start:          0,
		Length:         10000 * SectorSize,
		MetadataDevice: "/dev/sda1",
		DataDevice:     "/dev/sdb1",
		DataBlockSize:  128 * SectorSize,
		LowWaterMark:  100,
		Features:       []string{"skip_block_zeroing", "ignore_discard"},
	}
	require.Equal(t, "/dev/sda1 /dev/sdb1 128 100 2 skip_block_zeroing ignore_discard", tp.buildSpec())
}

func TestThinBuildSpec(t *testing.T) {
	t.Parallel()

	th := ThinTable{
		Start:      0,
		Length:     1000 * SectorSize,
		PoolDevice: "/dev/mapper/pool",
		DeviceID:   42,
	}
	require.Equal(t, "/dev/mapper/pool 42", th.buildSpec())
	require.Equal(t, "thin", th.targetType())
}

func TestThinBuildSpecWithExternalOrigin(t *testing.T) {
	t.Parallel()

	th := ThinTable{
		Start:          0,
		Length:         1000 * SectorSize,
		PoolDevice:     "/dev/mapper/pool",
		DeviceID:       7,
		ExternalOrigin: "/dev/sdc1",
	}
	require.Equal(t, "/dev/mapper/pool 7 /dev/sdc1", th.buildSpec())
}
