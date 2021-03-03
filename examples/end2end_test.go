package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/anatol/devmapper.go"
)

func devInfo(name string) (map[string]string, error) {
	out, err := exec.Command("dmsetup", "info", name).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("dminfo(%s): %v\n%s", name, err, out)
	}

	re := regexp.MustCompile(`([\w ,]+):\W*([^\n]+)`)
	matches := re.FindAllStringSubmatch(string(out), -1)

	result := make(map[string]string)
	for _, m := range matches {
		result[m[1]] = m[2]
	}

	return result, nil
}

const (
	PropName          = "Name"
	PropState         = "State"
	PropReadAhead     = "Read Ahead"
	PropTablesPresent = "Tables present"
	PropOpenCount     = "Open count"
	PropEventNumber   = "Event number"
	PropDevId         = "Major, minor"
	PropTargetsNum    = "Number of targets"
	PropUUID          = "UUID"
)

func checkDevInfo(t *testing.T, got, expected map[string]string) {
	for k, v := range expected {
		if got[k] != v {
			t.Fatalf("Property '%s': expected %s, got %s", k, v, got[k])
		}
	}
}

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

func TestZeroTarget(t *testing.T) {
	name := "test.zerotarget"
	uuid := "2f144136-b0de-4b51-b2eb-bd869cc39a6e"
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
