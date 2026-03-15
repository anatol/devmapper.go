package devmapper

import (
	"io/fs"
	"strconv"
	"strings"
)

// ThinPoolTable represents information needed for 'thin-pool' target creation.
// It manages a pool of blocks shared between thin devices.
type ThinPoolTable struct {
	Start         uint64
	Length        uint64
	MetadataDevice string
	DataDevice    string
	DataBlockSize uint64 // block size in bytes
	LowWaterMark uint64 // free space threshold in blocks
	Features      []string
}

func (t ThinPoolTable) start() uint64 {
	return t.Start
}

func (t ThinPoolTable) length() uint64 {
	return t.Length
}

func (t ThinPoolTable) targetType() string {
	return "thin-pool"
}

func (t ThinPoolTable) buildSpec() string {
	args := []string{
		t.MetadataDevice,
		t.DataDevice,
		strconv.FormatUint(t.DataBlockSize/SectorSize, 10),
		strconv.FormatUint(t.LowWaterMark, 10),
	}

	if len(t.Features) > 0 {
		args = append(args, strconv.Itoa(len(t.Features)))
		args = append(args, t.Features...)
	} else {
		args = append(args, "0")
	}

	return strings.Join(args, " ")
}

type thinPoolVolume struct{}

func (t ThinPoolTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &thinPoolVolume{}, nil
}

func (t thinPoolVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (t thinPoolVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (t thinPoolVolume) Close() error {
	return errNotImplemented
}

// ThinTable represents information needed for 'thin' target creation.
// It creates a thin provisioned device backed by a thin-pool.
type ThinTable struct {
	Start          uint64
	Length         uint64
	PoolDevice     string
	DeviceID       uint64
	ExternalOrigin string // optional external origin device
}

func (t ThinTable) start() uint64 {
	return t.Start
}

func (t ThinTable) length() uint64 {
	return t.Length
}

func (t ThinTable) targetType() string {
	return "thin"
}

func (t ThinTable) buildSpec() string {
	args := []string{
		t.PoolDevice,
		strconv.FormatUint(t.DeviceID, 10),
	}

	if t.ExternalOrigin != "" {
		args = append(args, t.ExternalOrigin)
	}

	return strings.Join(args, " ")
}

type thinVolume struct{}

func (t ThinTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &thinVolume{}, nil
}

func (t thinVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (t thinVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (t thinVolume) Close() error {
	return errNotImplemented
}
