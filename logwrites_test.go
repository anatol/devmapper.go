package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogWritesBuildSpec(t *testing.T) {
	t.Parallel()

	l := LogWritesTable{
		Start:     0,
		Length:    10000 * SectorSize,
		Device:    "/dev/sda1",
		LogDevice: "/dev/sdb1",
	}
	require.Equal(t, "/dev/sda1 /dev/sdb1", l.buildSpec())
	require.Equal(t, "log-writes", l.targetType())
	require.Equal(t, uint64(0), l.start())
	require.Equal(t, uint64(10000*SectorSize), l.length())
}
