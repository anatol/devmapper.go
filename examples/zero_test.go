package main

import (
	"testing"

	"github.com/anatol/devmapper.go"
)

func TestZeroTarget(t *testing.T) {
	name := "test.zerotarget"
	uuid := "2fa44836-b0de-4b51-b2eb-bd811cc39a6e"
	z := devmapper.ZeroTable{Length: 200}
	if err := devmapper.CreateAndLoad(name, uuid, z); err != nil {
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
