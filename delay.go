package devmapper

import (
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// DelayTable represents information needed for 'delay' target creation.
// It delays I/O for testing purposes.
type DelayTable struct {
	Start       uint64
	Length      uint64
	ReadDevice  string
	ReadOffset  uint64
	ReadDelay   uint64 // delay in milliseconds
	WriteDevice string // optional, defaults to ReadDevice
	WriteOffset uint64
	WriteDelay  uint64 // delay in milliseconds
	FlushDevice string // optional, defaults to WriteDevice
	FlushOffset uint64
	FlushDelay  uint64 // delay in milliseconds
}

func (d DelayTable) start() uint64 {
	return d.Start
}

func (d DelayTable) length() uint64 {
	return d.Length
}

func (d DelayTable) targetType() string {
	return "delay"
}

func (d DelayTable) buildSpec() string {
	args := []string{
		d.ReadDevice,
		strconv.FormatUint(d.ReadOffset/SectorSize, 10),
		strconv.FormatUint(d.ReadDelay, 10),
	}

	if d.WriteDevice != "" {
		args = append(args,
			d.WriteDevice,
			strconv.FormatUint(d.WriteOffset/SectorSize, 10),
			strconv.FormatUint(d.WriteDelay, 10),
		)

		if d.FlushDevice != "" {
			args = append(args,
				d.FlushDevice,
				strconv.FormatUint(d.FlushOffset/SectorSize, 10),
				strconv.FormatUint(d.FlushDelay, 10),
			)
		}
	}

	return strings.Join(args, " ")
}

type delayVolume struct {
	readFile  *os.File
	readOff   int64
	writeFile *os.File
	writeOff  int64
}

func (d DelayTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	readFile, err := os.OpenFile(d.ReadDevice, flag, perm)
	if err != nil {
		return nil, err
	}

	writeDevice := d.WriteDevice
	writeOffset := d.WriteOffset
	if writeDevice == "" {
		writeDevice = d.ReadDevice
		writeOffset = d.ReadOffset
	}

	var writeFile *os.File
	if writeDevice == d.ReadDevice && writeOffset == d.ReadOffset {
		writeFile = readFile
	} else {
		writeFile, err = os.OpenFile(writeDevice, flag, perm)
		if err != nil {
			readFile.Close()
			return nil, err
		}
	}

	return &delayVolume{
		readFile:  readFile,
		readOff:   int64(d.ReadOffset),
		writeFile: writeFile,
		writeOff:  int64(writeOffset),
	}, nil
}

func (d delayVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return d.readFile.ReadAt(p, off+d.readOff)
}

func (d delayVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return d.writeFile.WriteAt(p, off+d.writeOff)
}

func (d delayVolume) Close() error {
	err := d.readFile.Close()
	if d.writeFile != d.readFile {
		if err2 := d.writeFile.Close(); err == nil {
			err = err2
		}
	}
	return err
}
