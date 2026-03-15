package devmapper

import (
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// DustTable represents information needed for 'dust' target creation.
// It simulates failing sectors on a device for testing purposes.
type DustTable struct {
	Start     uint64
	Length    uint64
	Device    string
	Offset    uint64
	BlockSize uint64 // block size in bytes
}

func (d DustTable) start() uint64 {
	return d.Start
}

func (d DustTable) length() uint64 {
	return d.Length
}

func (d DustTable) targetType() string {
	return "dust"
}

func (d DustTable) buildSpec() string {
	args := []string{
		d.Device,
		strconv.FormatUint(d.Offset/SectorSize, 10),
		strconv.FormatUint(d.BlockSize, 10),
	}
	return strings.Join(args, " ")
}

type dustVolume struct {
	f      *os.File
	offset int64
}

func (d DustTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	f, err := os.OpenFile(d.Device, flag, perm)
	if err != nil {
		return nil, err
	}
	return &dustVolume{f: f, offset: int64(d.Offset)}, nil
}

func (d dustVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return d.f.ReadAt(p, off+d.offset)
}

func (d dustVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return d.f.WriteAt(p, off+d.offset)
}

func (d dustVolume) Close() error {
	return d.f.Close()
}
