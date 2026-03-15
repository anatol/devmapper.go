package devmapper

import (
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// StripeDevice represents one device in a stripe set
type StripeDevice struct {
	Device string
	Offset uint64
}

// StripeTable represents information needed for 'striped' target creation.
// It stripes I/O across multiple devices (RAID-0).
type StripeTable struct {
	Start     uint64
	Length    uint64
	ChunkSize uint64 // chunk size in bytes
	Devices   []StripeDevice
}

func (s StripeTable) start() uint64 {
	return s.Start
}

func (s StripeTable) length() uint64 {
	return s.Length
}

func (s StripeTable) targetType() string {
	return "striped"
}

func (s StripeTable) buildSpec() string {
	args := []string{
		strconv.Itoa(len(s.Devices)),
		strconv.FormatUint(s.ChunkSize/SectorSize, 10),
	}
	for _, d := range s.Devices {
		args = append(args, d.Device, strconv.FormatUint(d.Offset/SectorSize, 10))
	}
	return strings.Join(args, " ")
}

type stripeVolume struct {
	files     []*os.File
	offsets   []int64
	chunkSize uint64
}

func (s StripeTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	files := make([]*os.File, len(s.Devices))
	offsets := make([]int64, len(s.Devices))
	for i, d := range s.Devices {
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
	return &stripeVolume{files: files, offsets: offsets, chunkSize: s.ChunkSize}, nil
}

func (s stripeVolume) ReadAt(p []byte, off int64) (int, error) {
	return s.doIO(p, off, true)
}

func (s stripeVolume) WriteAt(p []byte, off int64) (int, error) {
	return s.doIO(p, off, false)
}

func (s stripeVolume) doIO(p []byte, off int64, read bool) (int, error) {
	n := len(s.files)
	chunk := s.chunkSize
	done := 0

	for done < len(p) {
		pos := uint64(off) + uint64(done)
		// which chunk in the virtual device
		chunkIdx := pos / chunk
		// offset within that chunk
		chunkOff := pos % chunk
		// which stripe device
		devIdx := int(chunkIdx % uint64(n))
		// which chunk on that device
		devChunk := chunkIdx / uint64(n)
		// final offset on the device
		devOff := int64(devChunk*chunk+chunkOff) + s.offsets[devIdx]
		// how much we can do within this chunk
		remaining := chunk - chunkOff
		toProcess := min(remaining, uint64(len(p)-done))

		buf := p[done : done+int(toProcess)]
		var nn int
		var err error
		if read {
			nn, err = s.files[devIdx].ReadAt(buf, devOff)
		} else {
			nn, err = s.files[devIdx].WriteAt(buf, devOff)
		}
		done += nn
		if err != nil {
			return done, err
		}
	}

	return done, nil
}

func (s stripeVolume) Close() error {
	for _, f := range s.files {
		f.Close()
	}
	return nil
}
