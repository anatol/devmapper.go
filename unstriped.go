package devmapper

import (
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// UnstripedTable represents information needed for 'unstriped' target creation.
// It extracts a single stripe from a striped device.
type UnstripedTable struct {
	Start      uint64
	Length     uint64
	NumStripes uint64
	ChunkSize  uint64 // chunk size in bytes
	StripeNum  uint64 // which stripe to extract (0-indexed)
	Device     string
	Offset     uint64
}

func (u UnstripedTable) start() uint64 {
	return u.Start
}

func (u UnstripedTable) length() uint64 {
	return u.Length
}

func (u UnstripedTable) targetType() string {
	return "unstriped"
}

func (u UnstripedTable) buildSpec() string {
	args := []string{
		strconv.FormatUint(u.NumStripes, 10),
		strconv.FormatUint(u.ChunkSize/SectorSize, 10),
		strconv.FormatUint(u.StripeNum, 10),
		u.Device,
		strconv.FormatUint(u.Offset/SectorSize, 10),
	}
	return strings.Join(args, " ")
}

type unstripedVolume struct {
	f          *os.File
	offset     int64
	numStripes uint64
	chunkSize  uint64
	stripeNum  uint64
}

func (u UnstripedTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	f, err := os.OpenFile(u.Device, flag, perm)
	if err != nil {
		return nil, err
	}
	return &unstripedVolume{
		f:          f,
		offset:     int64(u.Offset),
		numStripes: u.NumStripes,
		chunkSize:  u.ChunkSize,
		stripeNum:  u.StripeNum,
	}, nil
}

func (u unstripedVolume) ReadAt(p []byte, off int64) (int, error) {
	return u.doIO(p, off, true)
}

func (u unstripedVolume) WriteAt(p []byte, off int64) (int, error) {
	return u.doIO(p, off, false)
}

func (u unstripedVolume) doIO(p []byte, off int64, read bool) (int, error) {
	chunk := u.chunkSize
	done := 0

	for done < len(p) {
		pos := uint64(off) + uint64(done)
		// which chunk group in the unstriped view
		chunkIdx := pos / chunk
		// offset within that chunk
		chunkOff := pos % chunk
		// map to the underlying striped device:
		// the full stripe group * numStripes * chunkSize + stripeNum * chunkSize + chunkOff
		devOff := int64(chunkIdx*u.numStripes*chunk+u.stripeNum*chunk+chunkOff) + u.offset

		remaining := chunk - chunkOff
		toProcess := min(remaining, uint64(len(p)-done))

		buf := p[done : done+int(toProcess)]
		var nn int
		var err error
		if read {
			nn, err = u.f.ReadAt(buf, devOff)
		} else {
			nn, err = u.f.WriteAt(buf, devOff)
		}
		done += nn
		if err != nil {
			return done, err
		}
	}

	return done, nil
}

func (u unstripedVolume) Close() error {
	return u.f.Close()
}
