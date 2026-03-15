package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnstripedBuildSpec(t *testing.T) {
	t.Parallel()

	u := UnstripedTable{
		Start:      0,
		Length:     5000 * SectorSize,
		NumStripes: 4,
		ChunkSize:  128 * SectorSize,
		StripeNum:  2,
		Device:     "/dev/sda1",
		Offset:     0,
	}
	require.Equal(t, "4 128 2 /dev/sda1 0", u.buildSpec())
	require.Equal(t, "unstriped", u.targetType())
	require.Equal(t, uint64(0), u.start())
	require.Equal(t, uint64(5000*SectorSize), u.length())
}
