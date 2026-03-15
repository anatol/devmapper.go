package devmapper

import (
	"io/fs"
	"os"
	"strings"
)

// LogWritesTable represents information needed for 'log-writes' target creation.
// It logs all write operations to a separate device for crash consistency testing.
type LogWritesTable struct {
	Start     uint64
	Length    uint64
	Device    string
	LogDevice string
}

func (l LogWritesTable) start() uint64 {
	return l.Start
}

func (l LogWritesTable) length() uint64 {
	return l.Length
}

func (l LogWritesTable) targetType() string {
	return "log-writes"
}

func (l LogWritesTable) buildSpec() string {
	args := []string{l.Device, l.LogDevice}
	return strings.Join(args, " ")
}

type logWritesVolume struct {
	f *os.File
}

func (l LogWritesTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	f, err := os.OpenFile(l.Device, flag, perm)
	if err != nil {
		return nil, err
	}
	return &logWritesVolume{f: f}, nil
}

func (l logWritesVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return l.f.ReadAt(p, off)
}

func (l logWritesVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return l.f.WriteAt(p, off)
}

func (l logWritesVolume) Close() error {
	return l.f.Close()
}
