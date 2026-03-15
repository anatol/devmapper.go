package devmapper

import (
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// SnapshotOriginTable represents information needed for 'snapshot-origin' target creation.
// It marks a device as the origin for snapshots.
type SnapshotOriginTable struct {
	Start        uint64
	Length       uint64
	OriginDevice string
}

func (s SnapshotOriginTable) start() uint64 {
	return s.Start
}

func (s SnapshotOriginTable) length() uint64 {
	return s.Length
}

func (s SnapshotOriginTable) targetType() string {
	return "snapshot-origin"
}

func (s SnapshotOriginTable) buildSpec() string {
	return s.OriginDevice
}

type snapshotOriginVolume struct {
	f *os.File
}

func (s SnapshotOriginTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	f, err := os.OpenFile(s.OriginDevice, flag, perm)
	if err != nil {
		return nil, err
	}
	return &snapshotOriginVolume{f: f}, nil
}

func (s snapshotOriginVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return s.f.ReadAt(p, off)
}

func (s snapshotOriginVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return s.f.WriteAt(p, off)
}

func (s snapshotOriginVolume) Close() error {
	return s.f.Close()
}

// SnapshotTable represents information needed for 'snapshot' target creation.
// It creates a copy-on-write snapshot of the origin device.
type SnapshotTable struct {
	Start        uint64
	Length       uint64
	OriginDevice string
	COWDevice    string // copy-on-write device
	Persistent   bool   // persistent (P) or transient (N)
	ChunkSize    uint64 // chunk size in bytes
}

func (s SnapshotTable) start() uint64 {
	return s.Start
}

func (s SnapshotTable) length() uint64 {
	return s.Length
}

func (s SnapshotTable) targetType() string {
	return "snapshot"
}

func (s SnapshotTable) buildSpec() string {
	persistent := "N"
	if s.Persistent {
		persistent = "P"
	}
	args := []string{
		s.OriginDevice,
		s.COWDevice,
		persistent,
		strconv.FormatUint(s.ChunkSize/SectorSize, 10),
	}
	return strings.Join(args, " ")
}

type snapshotVolume struct{}

func (s SnapshotTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &snapshotVolume{}, nil
}

func (s snapshotVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (s snapshotVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (s snapshotVolume) Close() error {
	return errNotImplemented
}

// SnapshotMergeTable represents information needed for 'snapshot-merge' target creation.
// It merges a snapshot back into its origin.
type SnapshotMergeTable struct {
	Start        uint64
	Length       uint64
	OriginDevice string
	COWDevice    string
	Persistent   bool
	ChunkSize    uint64 // chunk size in bytes
}

func (s SnapshotMergeTable) start() uint64 {
	return s.Start
}

func (s SnapshotMergeTable) length() uint64 {
	return s.Length
}

func (s SnapshotMergeTable) targetType() string {
	return "snapshot-merge"
}

func (s SnapshotMergeTable) buildSpec() string {
	persistent := "N"
	if s.Persistent {
		persistent = "P"
	}
	args := []string{
		s.OriginDevice,
		s.COWDevice,
		persistent,
		strconv.FormatUint(s.ChunkSize/SectorSize, 10),
	}
	return strings.Join(args, " ")
}

type snapshotMergeVolume struct{}

func (s SnapshotMergeTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &snapshotMergeVolume{}, nil
}

func (s snapshotMergeVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (s snapshotMergeVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (s snapshotMergeVolume) Close() error {
	return errNotImplemented
}
