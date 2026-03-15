package devmapper

import (
	"io/fs"
	"strconv"
	"strings"
)

// CacheTable represents information needed for 'cache' target creation.
// It uses a fast device (e.g. SSD) to cache a slower origin device.
type CacheTable struct {
	Start          uint64
	Length         uint64
	MetadataDevice string
	CacheDevice    string
	OriginDevice   string
	BlockSize      uint64 // cache block size in bytes
	Features       []string
	Policy         string   // cache policy, e.g. "smq"
	PolicyArgs     []string // policy-specific arguments
}

func (c CacheTable) start() uint64 {
	return c.Start
}

func (c CacheTable) length() uint64 {
	return c.Length
}

func (c CacheTable) targetType() string {
	return "cache"
}

func (c CacheTable) buildSpec() string {
	args := []string{
		c.MetadataDevice,
		c.CacheDevice,
		c.OriginDevice,
		strconv.FormatUint(c.BlockSize/SectorSize, 10),
		strconv.Itoa(len(c.Features)),
	}
	args = append(args, c.Features...)
	args = append(args, c.Policy)
	args = append(args, strconv.Itoa(len(c.PolicyArgs)))
	args = append(args, c.PolicyArgs...)
	return strings.Join(args, " ")
}

type cacheVolume struct{}

func (c CacheTable) openVolume(flag int, perm fs.FileMode) (Volume, error) {
	return &cacheVolume{}, nil
}

func (c cacheVolume) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (c cacheVolume) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errNotImplemented
}

func (c cacheVolume) Close() error {
	return errNotImplemented
}
