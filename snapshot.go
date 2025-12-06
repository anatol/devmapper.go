package devmapper

import (
	"io/fs"
	"strconv"
	"strings"
)

type SnapshotOriginTable struct {
	Device string
	Start  uint64 // offset in bytes from which the snapshot origin starts
	Length uint64 // length in bytes of the snapshot origin
}

func (s SnapshotOriginTable) start() uint64      { return s.Start }
func (s SnapshotOriginTable) length() uint64     { return s.Length }
func (s SnapshotOriginTable) targetType() string { return "snapshot-origin" }
func (s SnapshotOriginTable) buildSpec() string {
	// <dev_path>
	return s.Device
}
func (s SnapshotOriginTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return nil, errNotImplemented
}

type SnapshotTable struct {
	OriginDevice string
	CowDevice    string
	Start        uint64
	Length       uint64

	Persistent bool
	Overflow   bool   // Persistent snapshots only: userspace supports overflow snapshots
	ChunkSize  uint64 // chunk size in sectors

	// Features
	DiscardZeroesCow      bool
	DiscardPassdownOrigin bool
}

func (s SnapshotTable) start() uint64      { return s.Start }
func (s SnapshotTable) length() uint64     { return s.Length }
func (s SnapshotTable) targetType() string { return "snapshot" }
func (s SnapshotTable) buildSpec() string {
	// <origin_dev> <COW-dev> <p|po|n> <chunk-size> [<# feature args> [<arg>]*]
	args := []string{s.OriginDevice, s.CowDevice}
	if s.Persistent {
		persistType := "P"
		if s.Overflow {
			persistType += "O"
		}
		args = append(args, persistType)
	} else {
		args = append(args, "N")
	}
	args = append(args, strconv.FormatUint(s.ChunkSize, 10))

	features := []string{}
	if s.DiscardZeroesCow {
		features = append(features, "discard_zeroes_cow")
	}
	if s.DiscardPassdownOrigin {
		features = append(features, "discard_passdown_origin")
	}
	args = append(args, strconv.Itoa(len(features)))
	args = append(args, features...)

	return strings.Join(args, " ")
}
func (s SnapshotTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return nil, errNotImplemented
}

type SnapshotMergeTable struct {
	SnapshotTable
}

func (s SnapshotMergeTable) start() uint64      { return s.Start }
func (s SnapshotMergeTable) length() uint64     { return s.Length }
func (s SnapshotMergeTable) targetType() string { return "snapshot-merge" }
func (s SnapshotMergeTable) buildSpec() string  { return s.SnapshotTable.buildSpec() }
func (s SnapshotMergeTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return nil, errNotImplemented
}
