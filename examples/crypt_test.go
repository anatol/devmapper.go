package main

import (
	"testing"

	"github.com/anatol/devmapper.go"
)

func TestCryptTarget(t *testing.T) {
	name := "test.crypttarget"
	uuid := "2f144136-b0de-4b51-b2eb-bd869cc39a6e"
	c := devmapper.CryptTable{
		Length:        600,
		Encryption:    "aes-xts-plain64",
		Key:           "babebabebabebabebabebabebabebabebabebabebabebabebabebabebabebabe",
		BackendDevice: "/dev/loop0",
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
