package devmapper

import (
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// FlakeyTable represents information needed for 'flakey' target creation.
// It creates an unreliable device for testing purposes, periodically
// returning errors for all I/O.
type FlakeyTable struct {
	Start        uint64
	Length       uint64
	Device       string
	Offset       uint64
	UpInterval   uint64 // number of seconds device is available
	DownInterval uint64 // number of seconds device returns errors
	Features     []string
}

func (f FlakeyTable) start() uint64 {
	return f.Start
}

func (f FlakeyTable) length() uint64 {
	return f.Length
}

func (f FlakeyTable) targetType() string {
	return "flakey"
}

func (f FlakeyTable) buildSpec() string {
	args := []string{
		f.Device,
		strconv.FormatUint(f.Offset/SectorSize, 10),
		strconv.FormatUint(f.UpInterval, 10),
		strconv.FormatUint(f.DownInterval, 10),
	}

	if len(f.Features) > 0 {
		args = append(args, strconv.Itoa(len(f.Features)))
		args = append(args, f.Features...)
	}

	return strings.Join(args, " ")
}

type flakeyVolume struct {
	f      *os.File
	offset int64
}

func (f FlakeyTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	file, err := os.OpenFile(f.Device, flag, perm)
	if err != nil {
		return nil, err
	}
	return &flakeyVolume{f: file, offset: int64(f.Offset)}, nil
}

func (f flakeyVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return f.f.ReadAt(p, off+f.offset)
}

func (f flakeyVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return f.f.WriteAt(p, off+f.offset)
}

func (f flakeyVolume) Close() error {
	return f.f.Close()
}
