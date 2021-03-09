package devmapper

import (
	"strconv"
	"strings"
)

const (
	// flags for crypt target
	CryptFlagAllowDiscards       = "allow_discards"
	CryptFlagSameCPUCrypt        = "same_cpu_crypt"
	CryptFlagSubmitFromCryptCPUs = "submit_from_crypt_cpus"
	CryptFlagNoReadWorkqueue     = "no_read_workqueue"
	CryptFlagNoWriteWorkqueue    = "no_write_workqueue"
)

// CryptTable represents information needed for 'crypt' target creation
type CryptTable struct {
	StartSector, Length uint64
	BackendDevice       string // device that stores the encrypted data
	BackendOffset       uint64
	Encryption          string
	Key                 string // it could be a plain key or keyID in the keystore as ":32:logon:foobarkey"
	IVTweak             uint64
	Flags               []string
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
	args := []string{c.Encryption, c.Key, strconv.FormatUint(c.IVTweak, 10), c.BackendDevice, strconv.FormatUint(c.BackendOffset, 10)}
	args = append(args, strconv.Itoa(len(c.Flags)))
	args = append(args, c.Flags...)

	return strings.Join(args, " ")
}
