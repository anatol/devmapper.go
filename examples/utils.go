package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"testing"
)

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
