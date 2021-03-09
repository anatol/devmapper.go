package devmapper

import (
	"strconv"
	"strings"
)

// LinearTable represents information needed for 'linear' target creation
type LinearTable struct {
	StartSector, Length uint64
	BackendDevice       string
	BackendOffset       uint64
}

func (l LinearTable) startSector() uint64 {
	return l.StartSector
}

func (l LinearTable) length() uint64 {
	return l.Length
}

func (l LinearTable) targetType() string {
	return "linear"
}

func (l LinearTable) buildSpec() string {
	args := []string{l.BackendDevice, strconv.FormatUint(l.BackendOffset, 10)}
	return strings.Join(args, " ")
}
