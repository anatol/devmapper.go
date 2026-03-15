package devmapper

import (
	"io/fs"
	"strconv"
	"strings"
)

// SwitchDevice represents one path in a switch target
type SwitchDevice struct {
	Device string
	Offset uint64
}

// SwitchTable represents information needed for 'switch' target creation.
// It allows switching between multiple backend devices on a per-region basis.
type SwitchTable struct {
	Start      uint64
	Length     uint64
	RegionSize uint64 // region size in bytes
	Devices    []SwitchDevice
}

func (s SwitchTable) start() uint64 {
	return s.Start
}

func (s SwitchTable) length() uint64 {
	return s.Length
}

func (s SwitchTable) targetType() string {
	return "switch"
}

func (s SwitchTable) buildSpec() string {
	args := []string{
		strconv.Itoa(len(s.Devices)),
		strconv.FormatUint(s.RegionSize/SectorSize, 10),
	}
	for _, d := range s.Devices {
		args = append(args, d.Device, strconv.FormatUint(d.Offset/SectorSize, 10))
	}
	return strings.Join(args, " ")
}

type switchVolume struct{}

func (s SwitchTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &switchVolume{}, nil
}

func (s switchVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (s switchVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (s switchVolume) Close() error {
	return errNotImplemented
}
