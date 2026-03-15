package devmapper

import (
	"io/fs"
	"strconv"
	"strings"
)

// EraTable represents information needed for 'era' target creation.
// It tracks which blocks were written within a given period (era).
type EraTable struct {
	Start          uint64
	Length         uint64
	MetadataDevice string
	OriginDevice   string
	BlockSize      uint64 // block size in bytes
}

func (e EraTable) start() uint64 {
	return e.Start
}

func (e EraTable) length() uint64 {
	return e.Length
}

func (e EraTable) targetType() string {
	return "era"
}

func (e EraTable) buildSpec() string {
	args := []string{
		e.MetadataDevice,
		e.OriginDevice,
		strconv.FormatUint(e.BlockSize/SectorSize, 10),
	}
	return strings.Join(args, " ")
}

type eraVolume struct{}

func (e EraTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &eraVolume{}, nil
}

func (e eraVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (e eraVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (e eraVolume) Close() error {
	return errNotImplemented
}
