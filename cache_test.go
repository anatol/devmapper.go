package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCacheBuildSpec(t *testing.T) {
	t.Parallel()

	c := CacheTable{
		Start:          0,
		Length:         10000 * SectorSize,
		MetadataDevice: "/dev/sda1",
		CacheDevice:    "/dev/sdb1",
		OriginDevice:   "/dev/sdc1",
		BlockSize:      256 * SectorSize,
		Policy:         "smq",
	}
	require.Equal(t, "/dev/sda1 /dev/sdb1 /dev/sdc1 256 0 smq 0", c.buildSpec())
	require.Equal(t, "cache", c.targetType())
	require.Equal(t, uint64(0), c.start())
	require.Equal(t, uint64(10000*SectorSize), c.length())
}

func TestCacheBuildSpecWithFeatures(t *testing.T) {
	t.Parallel()

	c := CacheTable{
		Start:          0,
		Length:         10000 * SectorSize,
		MetadataDevice: "/dev/sda1",
		CacheDevice:    "/dev/sdb1",
		OriginDevice:   "/dev/sdc1",
		BlockSize:      256 * SectorSize,
		Features:       []string{"writethrough"},
		Policy:         "smq",
		PolicyArgs:     []string{"migration_threshold", "2048"},
	}
	require.Equal(t, "/dev/sda1 /dev/sdb1 /dev/sdc1 256 1 writethrough smq 2 migration_threshold 2048", c.buildSpec())
}
