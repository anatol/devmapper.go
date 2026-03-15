package devmapper

import (
	"io/fs"
	"syscall"
)

// ErrorTable represents information needed for 'error' target creation.
// Any I/O to this target always fails with EIO.
type ErrorTable struct {
	Start  uint64
	Length uint64
}

func (e ErrorTable) start() uint64 {
	return e.Start
}

func (e ErrorTable) length() uint64 {
	return e.Length
}

func (e ErrorTable) targetType() string {
	return "error"
}

func (e ErrorTable) buildSpec() string {
	return ""
}

type errorVolume struct{}

func (e ErrorTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &errorVolume{}, nil
}

func (e errorVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, syscall.EIO
}

func (e errorVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, syscall.EIO
}

func (e errorVolume) Close() error {
	return nil
}
