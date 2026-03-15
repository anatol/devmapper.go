package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEraBuildSpec(t *testing.T) {
	t.Parallel()

	e := EraTable{
		Start:          0,
		Length:         10000 * SectorSize,
		MetadataDevice: "/dev/sda1",
		OriginDevice:   "/dev/sdb1",
		BlockSize:      256 * SectorSize,
	}
	require.Equal(t, "/dev/sda1 /dev/sdb1 256", e.buildSpec())
	require.Equal(t, "era", e.targetType())
	require.Equal(t, uint64(0), e.start())
	require.Equal(t, uint64(10000*SectorSize), e.length())
}
