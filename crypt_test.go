package devmapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCryptBuildSpec(t *testing.T) {
	t.Parallel()

	c := CryptTable{
		Start:         0,
		Length:        40 * SectorSize,
		BackendDevice: "/dev/loop0",
		BackendOffset: 0,
		Encryption:    "aes-xts-plain64",
		Key:           make([]byte, 32),
		IVTweak:       0,
	}
	spec := c.buildSpec()
	require.Equal(t, "aes-xts-plain64 0000000000000000000000000000000000000000000000000000000000000000 0 /dev/loop0 0 0", spec)
	require.Equal(t, "crypt", c.targetType())
}

func TestCryptBuildSpecWithFlags(t *testing.T) {
	t.Parallel()

	c := CryptTable{
		Length:        40 * SectorSize,
		BackendDevice: "/dev/loop0",
		Encryption:    "aes-xts-plain64",
		Key:           make([]byte, 32),
		Flags:         []string{CryptFlagAllowDiscards, CryptFlagNoReadWorkqueue},
	}
	spec := c.buildSpec()
	require.Contains(t, spec, "2 allow_discards no_read_workqueue")
}

func TestCryptBuildSpecWithSectorSize(t *testing.T) {
	t.Parallel()

	c := CryptTable{
		Length:        40 * SectorSize,
		BackendDevice: "/dev/loop0",
		Encryption:    "aes-xts-plain64",
		Key:           make([]byte, 32),
		SectorSize:    4096,
		Flags:         []string{CryptFlagAllowDiscards},
	}
	spec := c.buildSpec()
	require.Contains(t, spec, "2 allow_discards sector_size:4096")
}

func TestCryptBuildSpecDoesNotMutateFlags(t *testing.T) {
	t.Parallel()

	flags := make([]string, 1, 10) // extra capacity to trigger the bug
	flags[0] = CryptFlagAllowDiscards

	c := CryptTable{
		Length:        40 * SectorSize,
		BackendDevice: "/dev/loop0",
		Encryption:    "aes-xts-plain64",
		Key:           make([]byte, 32),
		SectorSize:    4096,
		Flags:         flags,
	}
	c.buildSpec()
	require.Equal(t, []string{CryptFlagAllowDiscards}, flags[:1])
	require.Equal(t, "", flags[1:2][0], "buildSpec must not mutate the original Flags backing array")
}

func TestCryptBuildSpecWithKeyID(t *testing.T) {
	t.Parallel()

	c := CryptTable{
		Length:        40 * SectorSize,
		BackendDevice: "/dev/loop0",
		Encryption:    "aes-xts-plain64",
		KeyID:         ":32:logon:foobarkey",
	}
	spec := c.buildSpec()
	require.Contains(t, spec, ":32:logon:foobarkey")
}
