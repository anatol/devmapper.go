package devmapper

import (
	"io/fs"
	"strconv"
	"strings"
)

// WritecacheTable represents information needed for 'writecache' target creation.
// It caches writes on a fast device and later flushes them to the slower origin device.
type WritecacheTable struct {
	Start       uint64
	Length      uint64
	Type        uint64 // 0 for pmem (DAX), 1 for SSD (block device)
	DataDevice  string
	CacheDevice string
	BlockSize   uint64 // block size in bytes
	Features    []string
}

func (w WritecacheTable) start() uint64 {
	return w.Start
}

func (w WritecacheTable) length() uint64 {
	return w.Length
}

func (w WritecacheTable) targetType() string {
	return "writecache"
}

func (w WritecacheTable) buildSpec() string {
	args := []string{
		strconv.FormatUint(w.Type, 10),
		w.DataDevice,
		w.CacheDevice,
		strconv.FormatUint(w.BlockSize/SectorSize, 10),
	}

	if len(w.Features) > 0 {
		args = append(args, strconv.Itoa(len(w.Features)))
		args = append(args, w.Features...)
	} else {
		args = append(args, "0")
	}

	return strings.Join(args, " ")
}

type writecacheVolume struct{}

func (w WritecacheTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &writecacheVolume{}, nil
}

func (w writecacheVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (w writecacheVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (w writecacheVolume) Close() error {
	return errNotImplemented
}
