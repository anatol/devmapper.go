package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZeroBuildSpec(t *testing.T) {
	t.Parallel()

	z := ZeroTable{Start: 0, Length: 200 * SectorSize}
	require.Equal(t, "", z.buildSpec())
	require.Equal(t, "zero", z.targetType())
	require.Equal(t, uint64(0), z.start())
	require.Equal(t, uint64(200*SectorSize), z.length())
}

func TestZeroVolumeReadAt(t *testing.T) {
	t.Parallel()

	z := ZeroTable{Length: 200 * SectorSize}
	v, err := z.openVolume(0, 0)
	require.NoError(t, err)
	defer v.Close()

	// fill buffer with non-zero data, verify ReadAt clears it
	buf := make([]byte, 512)
	for i := range buf {
		buf[i] = 0xff
	}
	n, err := v.ReadAt(buf, 0)
	require.NoError(t, err)
	require.Equal(t, 512, n)
	require.Equal(t, make([]byte, 512), buf)
}

func TestZeroVolumeWriteAt(t *testing.T) {
	t.Parallel()

	z := ZeroTable{Length: 200 * SectorSize}
	v, err := z.openVolume(0, 0)
	require.NoError(t, err)
	defer v.Close()

	buf := make([]byte, 512)
	n, err := v.WriteAt(buf, 0)
	require.NoError(t, err)
	require.Equal(t, 512, n)
}
