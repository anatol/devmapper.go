package devmapper

import (
	"strconv"
	"strings"
)

const (
	// flags for crypt target

	// CryptFlagAllowDiscards is an equivalent of 'allow_discards' crypt option
	CryptFlagAllowDiscards = "allow_discards"
	// CryptFlagSameCPUCrypt is an equivalent of 'same_cpu_crypt' crypt option
	CryptFlagSameCPUCrypt = "same_cpu_crypt"
	// CryptFlagSubmitFromCryptCPUs is an equivalent of 'submit_from_crypt_cpus' crypt option
	CryptFlagSubmitFromCryptCPUs = "submit_from_crypt_cpus"
	// CryptFlagNoReadWorkqueue is an equivalent of 'no_read_workqueue' crypt option
	CryptFlagNoReadWorkqueue = "no_read_workqueue"
	// CryptFlagNoWriteWorkqueue is an equivalent of 'no_write_workqueue' crypt option
	CryptFlagNoWriteWorkqueue = "no_write_workqueue"
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
