package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/anatol/devmapper.go"
	"github.com/tych0/go-losetup"
	"golang.org/x/sys/unix"
)

func TestVerity(t *testing.T) {
	dir := t.TempDir()

	// Setup data device
	d, err := os.Create(dir + "/data")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	// Write some data there
	if _, err := d.WriteString("Hello verity!!!!"); err != nil {
		t.Fatal(err)
	}

	dataSize := 4096 * devmapper.SectorSize
	if err := d.Truncate(int64(dataSize)); err != nil {
		t.Fatal(err)
	}
	dLoop, err := losetup.Attach(dir+"/data", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	defer dLoop.Detach()

	// Setup hash device
	h, err := os.Create(dir + "/hash")
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()
	if err := h.Truncate(100 * devmapper.SectorSize); err != nil {
		t.Fatal(err)
	}
	hLoop, err := losetup.Attach(dir+"/hash", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	defer hLoop.Detach()

	// format hash device
	salt, err := randomHex(32)
	if err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command("veritysetup", "format", "--no-superblock", "--salt", salt, dLoop.Path(), hLoop.Path()).Output()
	if err != nil {
		t.Fatal(err)
	}
	props, err := parseProperties(out)
	if err != nil {
		t.Fatal(err)
	}

	v := devmapper.VerityTable{
		Length:        4096,
		HashType:      1, // the same as props["Hash type"]
		DataDevice:    dLoop.Path(),
		HashDevice:    hLoop.Path(),
		DataBlockSize: 4096,     // 4096 is the default value, the same as props["Data block size"]
		HashBlockSize: 4096,     // 4096 is the default value, the same as props["Hash block size"]
		NumDataBlocks: 512,      // size of the device in "DataBlockSize" blocks
		Algorithm:     "sha256", // the same as props["Hash algorithm"]
		Salt:          salt,
		Digest:        props["Root hash"],
	}
	name := "test.verity"
	uuid := "79639465-407e-4af6-b76b-d98a70c0d5d3"
	if err := devmapper.CreateAndLoad(name, uuid, devmapper.ReadOnlyFlag, v); err != nil {
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
		PropState:         "ACTIVE (READ-ONLY)",
		PropTablesPresent: "LIVE",
		PropUUID:          uuid,
	})

	mapper := "/dev/mapper/" + name
	if err := waitForFile(mapper); err != nil {
		t.Fatal(err)
	}

	// Now verify the data read from the mapper
	data, err := os.ReadFile(mapper)
	if err != nil {
		t.Fatal(err)
	}
	expectedData := make([]byte, 4096*devmapper.SectorSize)
	copy(expectedData, "Hello verity!!!!")
	if bytes.Compare(expectedData, data) != 0 {
		t.Fatal("data read from the mapper differs from the backing file")
	}

	// Now corrupt the backing file (flip the first character from H to h)
	// verity should fail
	if _, err := d.WriteAt([]byte{'h'}, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := os.ReadFile(mapper); err == nil {
		t.Fatal("expected EIO if backing device is corrupted")
	} else {
		e := errors.Unwrap(err)
		if e != unix.EIO {
			t.Fatalf("unexpected error on verity corruption: %v", err)
		}
	}
}
