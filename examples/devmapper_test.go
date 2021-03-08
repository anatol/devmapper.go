package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/anatol/devmapper.go"
	"github.com/tych0/go-losetup"
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

// create 3 files with 5 sectors each (3*512 bytes) then join it
func TestJoinedDevices(t *testing.T) {
	name := "test.joinedtarget"
	uuid := "2c4c1a46-0c00-4430-b708-bdd16d34e7bf"

	dir := t.TempDir()

	targets := make([]devmapper.Table, 0, 3)

	for i := 0; i < 3; i++ {
		backingFile := dir + "/backing." + strconv.Itoa(i)
		f, err := os.Create(backingFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		text := fmt.Sprintf("Hello, world %d !!!", i)
		if _, err := f.WriteString(text); err != nil {
			t.Fatal(err)
		}
		if err := f.Truncate(5 * 512); err != nil {
			t.Fatal(err)
		}

		loop, err := losetup.Attach(backingFile, 0, false)
		if err != nil {
			t.Fatal(err)
		}
		defer loop.Detach()

		l := devmapper.LinearTable{
			StartSector:   5 * uint64(i),
			Length:        5,
			BackendDevice: loop.Path(),
		}

		targets = append(targets, l)
	}

	if err := devmapper.CreateAndLoad(name, uuid, targets...); err != nil {
		t.Fatal(err)
	}
	defer devmapper.Remove(name)

	got, err := devInfo(name)
	if err != nil {
		t.Fatal(err)
	}
	checkDevInfo(t, got, map[string]string{
		PropName:          name,
		PropTargetsNum:    "3",
		PropState:         "ACTIVE",
		PropTablesPresent: "LIVE",
		PropUUID:          uuid,
	})

	time.Sleep(time.Second)
	data, err := os.ReadFile("/dev/mapper/" + name)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 15*512 {
		t.Fatalf("invalid file length: expect %d, got %d", 15*512, len(data))
	}
	expectedData := make([]byte, 15*512)
	copy(expectedData[0*512:], "Hello, world 0 !!!")
	copy(expectedData[5*512:], "Hello, world 1 !!!")
	copy(expectedData[10*512:], "Hello, world 2 !!!")

	if bytes.Compare(expectedData, data) != 0 {
		t.Fatal("data read from the mapper differs from the backing file")
	}
}

// split backing device into 2 devmap devices
func TestSplitDevices(t *testing.T) {
	dir := t.TempDir()

	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := f.Truncate(50 * 512); err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteAt([]byte("Hello, world 0 !!!"), 0); err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteAt([]byte("Hello, world 1 !!!"), 7*512); err != nil {
		t.Fatal(err)
	}

	loop, err := losetup.Attach(backingFile, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	defer loop.Detach()

	name1 := "test.splittarget1"
	uuid1 := "8fb63cb9-e92c-473e-bad2-3dffd17859e7"
	l1 := devmapper.LinearTable{
		Length:        5,
		BackendDevice: loop.Path(),
		BackendOffset: 0,
	}
	if err := devmapper.CreateAndLoad(name1, uuid1, l1); err != nil {
		t.Fatal(err)
	}
	defer devmapper.Remove(name1)

	name2 := "test.splittarget2"
	uuid2 := "790dd808-3b13-46e1-9d4a-42053580010c"
	l2 := devmapper.LinearTable{
		Length:        45,
		BackendDevice: loop.Path(),
		BackendOffset: 5,
	}
	if err := devmapper.CreateAndLoad(name2, uuid2, l2); err != nil {
		t.Fatal(err)
	}
	defer devmapper.Remove(name2)

	got1, err := devInfo(name1)
	if err != nil {
		t.Fatal(err)
	}
	checkDevInfo(t, got1, map[string]string{
		PropName:          name1,
		PropTargetsNum:    "1",
		PropState:         "ACTIVE",
		PropTablesPresent: "LIVE",
		PropUUID:          uuid1,
	})

	got2, err := devInfo(name2)
	if err != nil {
		t.Fatal(err)
	}
	checkDevInfo(t, got2, map[string]string{
		PropName:          name2,
		PropTargetsNum:    "1",
		PropState:         "ACTIVE",
		PropTablesPresent: "LIVE",
		PropUUID:          uuid2,
	})

	time.Sleep(time.Second)

	data1, err := os.ReadFile("/dev/mapper/" + name1)
	if err != nil {
		t.Fatal(err)
	}
	if len(data1) != 5*512 {
		t.Fatalf("invalid file length: expect %d, got %d", 5*512, len(data1))
	}
	expectedData := make([]byte, 5*512)
	copy(expectedData[0*512:], "Hello, world 0 !!!")
	if bytes.Compare(expectedData, data1) != 0 {
		t.Fatal("data read from the mapper differs from the backing file")
	}

	data2, err := os.ReadFile("/dev/mapper/" + name2)
	if err != nil {
		t.Fatal(err)
	}
	if len(data2) != 45*512 {
		t.Fatalf("invalid file length: expect %d, got %d", 45*512, len(data2))
	}
	expectedData2 := make([]byte, 45*512)
	copy(expectedData2[2*512:], "Hello, world 1 !!!")
	if bytes.Compare(expectedData2, data2) != 0 {
		t.Fatal("data read from the mapper differs from the backing file")
	}
}