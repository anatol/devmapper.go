package main

import (
	"strings"
	"testing"

	"github.com/anatol/devmapper.go"
)

func TestCreate(t *testing.T) {
	name := "test.create"
	uuid := "634d7039-f339-4bbb-860a-e8a8c865a729"
	if err := devmapper.Create(name, uuid); err != nil {
		t.Fatal(err)
	}

	// check the target exists
	got, err := devInfo(name)
	if err != nil {
		t.Fatal(err)
	}
	checkDevInfo(t, got, map[string]string{
		PropName:          name,
		PropState:         "ACTIVE",
		PropTablesPresent: "None",
		PropUUID:          uuid,
	})

	if err := devmapper.Remove(name); err != nil {
		t.Fatal(err)
	}

	// check the target got removed
	_, err = devInfo(name)
	if err == nil || !strings.Contains(err.Error(), "Device does not exist.") {
		t.Fatalf("dm device %s should be removed", name)
	}
}
