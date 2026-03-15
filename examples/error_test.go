package test

import (
	"os"
	"syscall"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/stretchr/testify/require"
)

func TestErrorTarget(t *testing.T) {
	name := "test.errortarget"
	uuid := "a1b2c3d4-1111-2222-3333-444455556666"
	e := devmapper.ErrorTable{Length: 200 * devmapper.SectorSize}
	require.NoError(t, devmapper.CreateAndLoad(name, uuid, 0, e))
	defer devmapper.Remove(name)

	got, err := devInfo(name)
	require.NoError(t, err)
	checkDevInfo(t, got, map[string]string{
		PropName:          name,
		PropTargetsNum:    "1",
		PropState:         "ACTIVE",
		PropTablesPresent: "LIVE",
		PropUUID:          uuid,
	})
}

func TestUserspaceErrorTarget(t *testing.T) {
	e := devmapper.ErrorTable{Length: 200 * devmapper.SectorSize}
	v, err := devmapper.OpenUserspaceVolume(os.O_RDONLY, 0, e)
	require.NoError(t, err)
	defer v.Close()

	buf := make([]byte, 512)
	_, err = v.ReadAt(buf, 0)
	require.ErrorIs(t, err, syscall.EIO)

	_, err = v.WriteAt(buf, 0)
	require.ErrorIs(t, err, syscall.EIO)
}
