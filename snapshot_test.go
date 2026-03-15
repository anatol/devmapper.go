package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnapshotOriginBuildSpec(t *testing.T) {
	t.Parallel()

	s := SnapshotOriginTable{
		Start:        0,
		Length:       1000 * SectorSize,
		OriginDevice: "/dev/sda1",
	}
	require.Equal(t, "/dev/sda1", s.buildSpec())
	require.Equal(t, "snapshot-origin", s.targetType())
	require.Equal(t, uint64(0), s.start())
	require.Equal(t, uint64(1000*SectorSize), s.length())
}

func TestSnapshotBuildSpecPersistent(t *testing.T) {
	t.Parallel()

	s := SnapshotTable{
		Start:        0,
		Length:       1000 * SectorSize,
		OriginDevice: "/dev/sda1",
		COWDevice:    "/dev/sdb1",
		Persistent:   true,
		ChunkSize:    128 * SectorSize,
	}
	require.Equal(t, "/dev/sda1 /dev/sdb1 P 128", s.buildSpec())
	require.Equal(t, "snapshot", s.targetType())
}

func TestSnapshotBuildSpecTransient(t *testing.T) {
	t.Parallel()

	s := SnapshotTable{
		Start:        0,
		Length:       1000 * SectorSize,
		OriginDevice: "/dev/sda1",
		COWDevice:    "/dev/sdb1",
		Persistent:   false,
		ChunkSize:    64 * SectorSize,
	}
	require.Equal(t, "/dev/sda1 /dev/sdb1 N 64", s.buildSpec())
}

func TestSnapshotMergeBuildSpec(t *testing.T) {
	t.Parallel()

	s := SnapshotMergeTable{
		Start:        0,
		Length:       1000 * SectorSize,
		OriginDevice: "/dev/sda1",
		COWDevice:    "/dev/sdb1",
		Persistent:   true,
		ChunkSize:    128 * SectorSize,
	}
	require.Equal(t, "/dev/sda1 /dev/sdb1 P 128", s.buildSpec())
	require.Equal(t, "snapshot-merge", s.targetType())
}
