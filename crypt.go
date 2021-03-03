package devmapper

import (
	"strconv"
	"strings"
)

// CryptTable represents information needed for 'crypt' target creation
type CryptTable struct {
	StartSector, Length uint64
	BackendDevice       string // device that stores the encrypted data
	BackendOffset       int
	Encryption          string
	Key                 string // it could be a plain key or keyID in the keystore as ":32:logon:foobarkey"
	IVTweak             int
	Flags               []string // Flags like allow_discards, same_cpu_crypt, ...
}

func (c CryptTable) startSector() uint64 {
	return c.StartSector
}

func (c CryptTable) length() uint64 {
	return c.Length
}

func (c CryptTable) targetType() string {
	return "crypt"
}

func (c CryptTable) buildSpec() string {
	storageArg := []string{c.Encryption, c.Key, strconv.Itoa(c.IVTweak), c.BackendDevice, strconv.Itoa(c.BackendOffset)}
	storageArg = append(storageArg, strconv.Itoa(len(c.Flags)))
	storageArg = append(storageArg, c.Flags...)

	return strings.Join(storageArg, " ")
}
