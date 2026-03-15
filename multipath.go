package devmapper

import (
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// MultipathPath represents a single path in a multipath path group
type MultipathPath struct {
	Device string
	Args   []string // per-path selector args
}

// MultipathGroup represents a path group in a multipath target
type MultipathGroup struct {
	PathSelector string // e.g. "round-robin"
	SelectorArgs []string
	Paths        []MultipathPath
}

// MultipathTable represents information needed for 'multipath' target creation.
// It provides I/O failover and load balancing across multiple paths.
type MultipathTable struct {
	Start      uint64
	Length     uint64
	Features   []string
	HWHandler  string   // hardware handler, e.g. "0" for none
	HWArgs     []string
	PathGroups []MultipathGroup
}

func (m MultipathTable) start() uint64 {
	return m.Start
}

func (m MultipathTable) length() uint64 {
	return m.Length
}

func (m MultipathTable) targetType() string {
	return "multipath"
}

func (m MultipathTable) buildSpec() string {
	args := []string{strconv.Itoa(len(m.Features))}
	args = append(args, m.Features...)
	args = append(args, m.HWHandler)
	args = append(args, strconv.Itoa(len(m.HWArgs)))
	args = append(args, m.HWArgs...)
	args = append(args, strconv.Itoa(len(m.PathGroups)))
	for _, g := range m.PathGroups {
		args = append(args, g.PathSelector, strconv.Itoa(len(g.SelectorArgs)))
		args = append(args, g.SelectorArgs...)
		args = append(args, strconv.Itoa(len(g.Paths)))
		for _, p := range g.Paths {
			args = append(args, p.Device, strconv.Itoa(len(p.Args)))
			args = append(args, p.Args...)
		}
	}
	return strings.Join(args, " ")
}

type multipathVolume struct {
	f *os.File
}

func (m MultipathTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	// Use the first path of the first path group
	if len(m.PathGroups) == 0 || len(m.PathGroups[0].Paths) == 0 {
		return nil, errNotImplemented
	}
	f, err := os.OpenFile(m.PathGroups[0].Paths[0].Device, flag, perm)
	if err != nil {
		return nil, err
	}
	return &multipathVolume{f: f}, nil
}

func (m multipathVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return m.f.ReadAt(p, off)
}

func (m multipathVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return m.f.WriteAt(p, off)
}

func (m multipathVolume) Close() error {
	return m.f.Close()
}
