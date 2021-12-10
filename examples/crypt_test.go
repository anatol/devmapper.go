package test

import (
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/stretchr/testify/require"
	"github.com/tych0/go-losetup"
)

func TestCryptTarget(t *testing.T) {
	name := "test.crypttarget"
	uuid := "2f144136-b0de-4b51-b2eb-bd869cc39a6e"

	dir := t.TempDir()
	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	require.NoError(t, err)
	defer f.Close()

	size := uint64(40) * devmapper.SectorSize
	require.NoError(t, f.Truncate(int64(size)))

	loop, err := losetup.Attach(backingFile, 0, false)
	require.NoError(t, err)
	defer loop.Detach()

	key := make([]byte, 32)
	c := devmapper.CryptTable{
		Length:        size,
		Encryption:    "aes-xts-plain64",
		Key:           key,
		BackendDevice: loop.Path(),
		Flags: []string{
			devmapper.CryptFlagAllowDiscards,
			devmapper.CryptFlagNoReadWorkqueue,
			devmapper.CryptFlagNoWriteWorkqueue,
			devmapper.CryptFlagSameCPUCrypt,
			devmapper.CryptFlagSubmitFromCryptCPUs,
		},
	}
	require.NoError(t, devmapper.CreateAndLoad(name, uuid, 0, c))
	defer devmapper.Remove(name)

	got, err := devInfo(name)
	require.NoError(t, err)
	checkDevInfo(t, got, map[string]string{
		PropName:          name,
		PropTargetsNum:    "1",
		PropState:         "ACTIVE",
		PropTablesPresent: "LIVE",
		PropUUID:          uuid,
	})
}
