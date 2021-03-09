package main

import (
	os "os"
	"testing"

	"github.com/anatol/devmapper.go"
)

func TestZeroTarget(t *testing.T) {
	name := "test.zerotarget"
	uuid := "2fa44836-b0de-4b51-b2eb-bd811cc39a6e"
	z := devmapper.ZeroTable{Length: 200}
	if err := devmapper.CreateAndLoad(name, uuid, 0, z); err != nil {
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

	mapper := "/dev/mapper/" + name
	if err := waitForFile(mapper); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(mapper)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) != 200*devmapper.SectorSize {
		t.Fatalf("expected size of the file %d, got %d", 200*devmapper.SectorSize, len(data))
	}
	for i, b := range data {
		if b != 0 {
			t.Fatalf("zero file must provide zeros, but got %d at index %d", b, i)
		}
	}
}
