package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinearBuildSpec(t *testing.T) {
	t.Parallel()

	l := LinearTable{
		Start:         0,
		Length:        10 * SectorSize,
		BackendDevice: "/dev/loop0",
		BackendOffset: 0,
	}
	require.Equal(t, "/dev/loop0 0", l.buildSpec())
	require.Equal(t, "linear", l.targetType())
	require.Equal(t, uint64(0), l.start())
	require.Equal(t, uint64(10*SectorSize), l.length())

	l2 := LinearTable{
		Start:         5 * SectorSize,
		Length:        20 * SectorSize,
		BackendDevice: "/dev/sda1",
		BackendOffset: 1024 * SectorSize,
	}
	require.Equal(t, "/dev/sda1 1024", l2.buildSpec())
	require.Equal(t, uint64(5*SectorSize), l2.start())
}
