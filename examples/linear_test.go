package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/tych0/go-losetup"
)

func TestLinear(t *testing.T) {
	name := "test.lineartarget"
	uuid := "b7d95b61-e077-4452-b904-13bd6936b44d"

	dir := t.TempDir()
	backingFile := dir + "/backing"
	f, err := os.Create(backingFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	text := "Hello, world!!!"
	if _, err := f.WriteString(text); err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(10 * 1024); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	loop, err := losetup.Attach(backingFile, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	defer loop.Detach()

	l := devmapper.LinearTable{
		Length:        15,
		BackendDevice: loop.Path(),
	}
	if err := devmapper.CreateAndLoad(name, uuid, l); err != nil {
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
	if len(data) != 15*devmapper.SectorSize {
		t.Fatalf("invalid file length: expect %d, got %d", 15*devmapper.SectorSize, len(data))
	}
	expectedData := make([]byte, 15*devmapper.SectorSize)
	copy(expectedData, text)

	if bytes.Compare(expectedData, data) != 0 {
		t.Fatal("data read from the mapper differs from the backing file")
	}
}
