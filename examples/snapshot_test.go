package test

import (
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/freddierice/go-losetup/v2"
	"github.com/stretchr/testify/require"
)

func TestSnapshotCreate(t *testing.T) {
	dir := t.TempDir()

	const originSize = 100 * devmapper.SectorSize
	originFile := dir + "/backing-origin"
	orig, err := os.Create(originFile)
	require.NoError(t, err)
	defer orig.Close()
	require.NoError(t, orig.Truncate(int64(originSize)))
	originLoop, err := losetup.Attach(originFile, 0, false)
	require.NoError(t, err)
	defer originLoop.Detach()

	const cowSize = 160 * devmapper.SectorSize
	cowFile := dir + "/backing-cow"
	cow, err := os.Create(cowFile)
	require.NoError(t, err)
	defer cow.Close()
	require.NoError(t, cow.Truncate(int64(cowSize)))
	cowLoop, err := losetup.Attach(cowFile, 0, false)
	require.NoError(t, err)
	defer cowLoop.Detach()

	originName := "test.origin"
	uuid := "2f141136-b0de-4b51-b2eb-bd849cc39a6e"
	o := devmapper.SnapshotOriginTable{
		Device: originLoop.Path(),
		Start:  0,
		Length: uint64(originSize),
	}
	require.NoError(t, devmapper.CreateAndLoad(originName, uuid, 0, o))
	defer devmapper.Remove(originName)

	gotOrigin, err := devInfo(originName)
	require.NoError(t, err)
	checkDevInfo(t, gotOrigin, map[string]string{
		PropName:          originName,
		PropTargetsNum:    "1",
		PropState:         "ACTIVE",
		PropTablesPresent: "LIVE",
		PropUUID:          uuid,
	})
	originDmPath := "/dev/mapper/" + originName
	require.NoError(t, waitForFile(originDmPath))

	// now create a snapshot and write something to origin
	snapshotName := "test.snapshot"
	snapshotUuid := "3f141136-b0de-4b51-b2eb-bd849cc39a6e"
	require.NoError(t, err)
	s := devmapper.SnapshotTable{
		OriginDevice: originDmPath,
		CowDevice:    cowLoop.Path(),
		Start:        0,
		Length:       uint64(originSize),
		Persistent:   true,
		ChunkSize:    4,
	}
	require.NoError(t, devmapper.CreateAndLoad(snapshotName, snapshotUuid, 0, s))
	defer devmapper.Remove(snapshotName)

	gotSnapshot, err := devInfo(snapshotName)
	require.NoError(t, err)
	checkDevInfo(t, gotSnapshot, map[string]string{
		PropName:          snapshotName,
		PropTargetsNum:    "1",
		PropState:         "ACTIVE",
		PropTablesPresent: "LIVE",
		PropUUID:          snapshotUuid,
	})
}
