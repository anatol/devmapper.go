package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorBuildSpec(t *testing.T) {
	t.Parallel()

	e := ErrorTable{Start: 0, Length: 200 * SectorSize}
	require.Equal(t, "", e.buildSpec())
	require.Equal(t, "error", e.targetType())
	require.Equal(t, uint64(0), e.start())
	require.Equal(t, uint64(200*SectorSize), e.length())
}

func TestErrorVolume(t *testing.T) {
	t.Parallel()

	e := ErrorTable{Length: 200 * SectorSize}
	v, err := e.openVolume(0, 0)
	require.NoError(t, err)
	defer v.Close()

	buf := make([]byte, 512)
	_, err = v.ReadAt(buf, 0)
	require.Error(t, err)

	_, err = v.WriteAt(buf, 0)
	require.Error(t, err)
}
