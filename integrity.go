package devmapper

import (
	"io/fs"
	"strconv"
	"strings"
)

// IntegrityTable represents information needed for 'integrity' target creation.
// It provides block-level data integrity protection (dm-integrity).
type IntegrityTable struct {
	Start    uint64
	Length   uint64
	Device   string
	Offset   uint64
	TagSize  uint64
	Mode     string // "J" (journal), "B" (bitmap), "D" (direct)
	Features []string
}

func (i IntegrityTable) start() uint64 {
	return i.Start
}

func (i IntegrityTable) length() uint64 {
	return i.Length
}

func (i IntegrityTable) targetType() string {
	return "integrity"
}

func (i IntegrityTable) buildSpec() string {
	args := []string{
		i.Device,
		strconv.FormatUint(i.Offset/SectorSize, 10),
		strconv.FormatUint(i.TagSize, 10),
		i.Mode,
	}

	if len(i.Features) > 0 {
		args = append(args, strconv.Itoa(len(i.Features)))
		args = append(args, i.Features...)
	} else {
		args = append(args, "0")
	}

	return strings.Join(args, " ")
}

type integrityVolume struct{}

func (i IntegrityTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &integrityVolume{}, nil
}

func (i integrityVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (i integrityVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (i integrityVolume) Close() error {
	return errNotImplemented
}
