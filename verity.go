package devmapper

import (
	"strconv"
	"strings"
)

// VerityTable represents information needed for 'verity' target creation
type VerityTable struct {
	StartSector, Length           uint64
	HashType                      uint64
	DataDevice                    string // the device containing data, the integrity of which needs to be checked
	HashDevice                    string // device that supplies the hash tree data
	DataBlockSize, HashBlockSize  uint64
	NumDataBlocks, HashStartBlock uint64
	Algorithm, Digest, Salt       string
	Params                        []string
}

func (v VerityTable) startSector() uint64 {
	return v.StartSector
}

func (v VerityTable) length() uint64 {
	return v.Length
}

func (v VerityTable) targetType() string {
	return "verity"
}

func (v VerityTable) buildSpec() string {
	args := []string{
		strconv.FormatUint(v.HashType, 10), v.DataDevice, v.HashDevice, strconv.FormatUint(v.DataBlockSize, 10),
		strconv.FormatUint(v.HashBlockSize, 10), strconv.FormatUint(v.NumDataBlocks, 10),
		strconv.FormatUint(v.HashStartBlock, 10), v.Algorithm, v.Digest, v.Salt,
	}

	args = append(args, v.Params...)
	return strings.Join(args, " ")
}
