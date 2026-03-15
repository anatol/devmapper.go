package devmapper

import (
	"io/fs"
	"strconv"
	"strings"
)

// RaidDevice represents one device pair in a RAID set (metadata + data)
type RaidDevice struct {
	MetaDevice string // metadata device, or "-" if not used
	DataDevice string
}

// RaidTable represents information needed for 'raid' target creation.
// It implements various RAID levels using the md (multiple devices) framework.
type RaidTable struct {
	Start    uint64
	Length   uint64
	RaidType string // e.g. "raid0", "raid1", "raid4", "raid5_ls", "raid6_zr", "raid10"
	Params   []string
	Devices  []RaidDevice
}

func (r RaidTable) start() uint64 {
	return r.Start
}

func (r RaidTable) length() uint64 {
	return r.Length
}

func (r RaidTable) targetType() string {
	return "raid"
}

func (r RaidTable) buildSpec() string {
	args := []string{
		r.RaidType,
		strconv.Itoa(len(r.Params)),
	}
	args = append(args, r.Params...)
	args = append(args, strconv.Itoa(len(r.Devices)))
	for _, d := range r.Devices {
		args = append(args, d.MetaDevice, d.DataDevice)
	}
	return strings.Join(args, " ")
}

type raidVolume struct{}

func (r RaidTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &raidVolume{}, nil
}

func (r raidVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (r raidVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (r raidVolume) Close() error {
	return errNotImplemented
}
