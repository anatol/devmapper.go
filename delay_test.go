package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDelayBuildSpecReadOnly(t *testing.T) {
	t.Parallel()

	d := DelayTable{
		Start:      0,
		Length:     1000 * SectorSize,
		ReadDevice: "/dev/sda1",
		ReadOffset: 0,
		ReadDelay:  100,
	}
	require.Equal(t, "/dev/sda1 0 100", d.buildSpec())
	require.Equal(t, "delay", d.targetType())
	require.Equal(t, uint64(0), d.start())
	require.Equal(t, uint64(1000*SectorSize), d.length())
}

func TestDelayBuildSpecReadWrite(t *testing.T) {
	t.Parallel()

	d := DelayTable{
		Start:       0,
		Length:      1000 * SectorSize,
		ReadDevice:  "/dev/sda1",
		ReadOffset:  0,
		ReadDelay:   100,
		WriteDevice: "/dev/sda1",
		WriteOffset: 0,
		WriteDelay:  200,
	}
	require.Equal(t, "/dev/sda1 0 100 /dev/sda1 0 200", d.buildSpec())
}

func TestDelayBuildSpecReadWriteFlush(t *testing.T) {
	t.Parallel()

	d := DelayTable{
		Start:       0,
		Length:      1000 * SectorSize,
		ReadDevice:  "/dev/sda1",
		ReadOffset:  0,
		ReadDelay:   50,
		WriteDevice: "/dev/sda1",
		WriteOffset: 0,
		WriteDelay:  100,
		FlushDevice: "/dev/sda1",
		FlushOffset: 0,
		FlushDelay:  150,
	}
	require.Equal(t, "/dev/sda1 0 50 /dev/sda1 0 100 /dev/sda1 0 150", d.buildSpec())
}
