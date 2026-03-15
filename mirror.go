package devmapper

import (
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// MirrorDevice represents one device in a mirror set
type MirrorDevice struct {
	Device string
	Offset uint64
}

// MirrorTable represents information needed for 'mirror' target creation.
// It mirrors I/O across multiple devices (RAID-1).
type MirrorTable struct {
	Start    uint64
	Length   uint64
	LogType  string   // e.g. "core" or "disk"
	LogArgs  []string // arguments for the log type
	Devices  []MirrorDevice
	Features []string // optional features like "handle_errors"
}

func (m MirrorTable) start() uint64 {
	return m.Start
}

func (m MirrorTable) length() uint64 {
	return m.Length
}

func (m MirrorTable) targetType() string {
	return "mirror"
}

func (m MirrorTable) buildSpec() string {
	args := []string{
		m.LogType,
		strconv.Itoa(len(m.LogArgs)),
	}
	args = append(args, m.LogArgs...)
	args = append(args, strconv.Itoa(len(m.Devices)))
	for _, d := range m.Devices {
		args = append(args, d.Device, strconv.FormatUint(d.Offset/SectorSize, 10))
	}
	args = append(args, strconv.Itoa(len(m.Features)))
	args = append(args, m.Features...)
	return strings.Join(args, " ")
}

type mirrorVolume struct {
	files   []*os.File
	offsets []int64
}

func (m MirrorTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	files := make([]*os.File, len(m.Devices))
	offsets := make([]int64, len(m.Devices))
	for i, d := range m.Devices {
		f, err := os.OpenFile(d.Device, flag, perm)
		if err != nil {
			for j := 0; j < i; j++ {
				files[j].Close()
			}
			return nil, err
		}
		files[i] = f
		offsets[i] = int64(d.Offset)
	}
	return &mirrorVolume{files: files, offsets: offsets}, nil
}

func (m mirrorVolume) ReadAt(p []byte, off int64) (n int, err error) {
	// Read from the first mirror device
	return m.files[0].ReadAt(p, off+m.offsets[0])
}

func (m mirrorVolume) WriteAt(p []byte, off int64) (n int, err error) {
	// Write to all mirror devices
	for i, f := range m.files {
		n, err = f.WriteAt(p, off+m.offsets[i])
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (m mirrorVolume) Close() error {
	var firstErr error
	for _, f := range m.files {
		if err := f.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
