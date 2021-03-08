package main

import (
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/tych0/go-losetup"
)

func TestCryptTarget(t *testing.T) {
	name := "test.crypttarget"
	uuid := "2f144136-b0de-4b51-b2eb-bd869cc39a6e"

	dir := t.TempDir()
	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := f.Truncate(40 * devmapper.SectorSize); err != nil {
		t.Fatal(err)
	}

	loop, err := losetup.Attach(backingFile, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	defer loop.Detach()

	c := devmapper.CryptTable{
		Length:        40,
		Encryption:    "aes-xts-plain64",
		Key:           "babebabebabebabebabebabebabebabebabebabebabebabebabebabebabebabe",
		BackendDevice: loop.Path(),
		Flags:         []string{"allow_discards"},
	}
	if err := devmapper.CreateAndLoad(name, uuid, c); err != nil {
		t.Fatal(err)
	}
	defer devmapper.Remove(name)

	got, err := devInfo(name)
	if err != nil {
		t.Fatal(err)
	}
	checkDevInfo(t, got, map[string]string{
		PropName:          name,
		PropTargetsNum:    "1",
		PropState:         "ACTIVE",
		PropTablesPresent: "LIVE",
		PropUUID:          uuid,
	})
}
