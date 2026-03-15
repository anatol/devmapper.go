package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerityBuildSpec(t *testing.T) {
	t.Parallel()

	v := VerityTable{
		Start:         0,
		Length:        4096 * SectorSize,
		HashType:      1,
		DataDevice:    "/dev/loop0",
		HashDevice:    "/dev/loop1",
		DataBlockSize: 4096,
		HashBlockSize: 4096,
		NumDataBlocks: 512,
		HashStartBlock: 0,
		Algorithm:     "sha256",
		Digest:        "abc123",
		Salt:          "def456",
	}
	expected := "1 /dev/loop0 /dev/loop1 4096 4096 512 0 sha256 abc123 def456"
	require.Equal(t, expected, v.buildSpec())
	require.Equal(t, "verity", v.targetType())
	require.Equal(t, uint64(0), v.start())
	require.Equal(t, uint64(4096*SectorSize), v.length())
}

func TestVerityBuildSpecWithParams(t *testing.T) {
	t.Parallel()

	v := VerityTable{
		Length:        4096 * SectorSize,
		HashType:      1,
		DataDevice:    "/dev/loop0",
		HashDevice:    "/dev/loop1",
		DataBlockSize: 4096,
		HashBlockSize: 4096,
		NumDataBlocks: 512,
		Algorithm:     "sha256",
		Digest:        "abc123",
		Salt:          "def456",
		Params:        []string{"1", "ignore_corruption"},
	}
	spec := v.buildSpec()
	require.Contains(t, spec, "def456 1 ignore_corruption")
}
