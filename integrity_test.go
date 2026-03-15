package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrityBuildSpec(t *testing.T) {
	t.Parallel()

	i := IntegrityTable{
		Start:   0,
		Length:  10000 * SectorSize,
		Device:  "/dev/sda1",
		Offset:  0,
		TagSize: 32,
		Mode:    "J",
	}
	require.Equal(t, "/dev/sda1 0 32 J 0", i.buildSpec())
	require.Equal(t, "integrity", i.targetType())
	require.Equal(t, uint64(0), i.start())
	require.Equal(t, uint64(10000*SectorSize), i.length())
}

func TestIntegrityBuildSpecWithFeatures(t *testing.T) {
	t.Parallel()

	i := IntegrityTable{
		Start:    0,
		Length:   10000 * SectorSize,
		Device:   "/dev/sda1",
		Offset:   0,
		TagSize:  32,
		Mode:     "B",
		Features: []string{"internal_hash:crc32c", "journal_sectors:2048"},
	}
	require.Equal(t, "/dev/sda1 0 32 B 2 internal_hash:crc32c journal_sectors:2048", i.buildSpec())
}
