package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWritecacheBuildSpec(t *testing.T) {
	t.Parallel()

	w := WritecacheTable{
		Start:       0,
		Length:      10000 * SectorSize,
		Type:        1,
		DataDevice:  "/dev/sda1",
		CacheDevice: "/dev/sdb1",
		BlockSize:   8 * SectorSize,
	}
	require.Equal(t, "1 /dev/sda1 /dev/sdb1 8 0", w.buildSpec())
	require.Equal(t, "writecache", w.targetType())
	require.Equal(t, uint64(0), w.start())
	require.Equal(t, uint64(10000*SectorSize), w.length())
}

func TestWritecacheBuildSpecWithFeatures(t *testing.T) {
	t.Parallel()

	w := WritecacheTable{
		Start:       0,
		Length:      10000 * SectorSize,
		Type:        0,
		DataDevice:  "/dev/sda1",
		CacheDevice: "/dev/pmem0",
		BlockSize:   8 * SectorSize,
		Features:    []string{"high_watermark", "50", "low_watermark", "25"},
	}
	require.Equal(t, "0 /dev/sda1 /dev/pmem0 8 4 high_watermark 50 low_watermark 25", w.buildSpec())
}
