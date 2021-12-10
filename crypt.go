package devmapper

import (
	"encoding/hex"
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
	Start         uint64
	Length        uint64
	BackendDevice string // device that stores the encrypted data
	BackendOffset uint64
	Encryption    string
	Key           []byte
	KeyID         string // key id in the keystore e.g. ":32:logon:foobarkey"
	IVTweak       uint64
	Flags         []string
}

func (c CryptTable) start() uint64 {
	return c.Start
}

func (c CryptTable) length() uint64 {
	return c.Length
}

func (c CryptTable) targetType() string {
	return "crypt"
}

func (c CryptTable) buildSpec() string {
	key := c.KeyID
	if key == "" {
		// dm-crypt requires hex-encoded password
		key = hex.EncodeToString(c.Key)
	}

	args := []string{c.Encryption, key, strconv.FormatUint(c.IVTweak, 10), c.BackendDevice, strconv.FormatUint(c.BackendOffset/SectorSize, 10)}
	args = append(args, strconv.Itoa(len(c.Flags)))
	args = append(args, c.Flags...)

	return strings.Join(args, " ")
}
