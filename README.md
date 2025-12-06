# Pure Go library for device mapper targets management

`devmapper.go` is a pure-Go library that helps to deal with device mapper targets.

Here is an example that demonstrates the API usage:
```go
func main() {
    name := "crypttarget"
    uuid := "2f144136-b0de-4b51-b2eb-bd869cc39a6e"
    key := make([]byte, 32)
    c := devmapper.CryptTable{
        Length:        60000 * 512, // size of the device in bytes
        Encryption:    "aes-xts-plain64",
        Key:           key,
        BackendDevice: "/dev/loop0",
        Flags:         []string{devmapper.CryptFlagAllowDiscards},
    }
    if err := devmapper.CreateAndLoad(name, uuid, c); err != nil {
        // handle error
    }
    defer devmapper.Remove(name)

    // at this point a devmapper target named 'crypttarget' should exist
    // you can check it with 'dmsetup info crypttarget'
    // and udev will create /dev/mapper/crypttarget device file
}
```

Or the same crypttarget initialization using Linux keychain

```go
func main() {
    // load key into keyring
    keyname := fmt.Sprintf("cryptsetup:%s-d%d", uuid, luksDigestId) // an example of keyname used by LUKS framework
    kid, err := unix.AddKey("logon", keyname, key, unix.KEY_SPEC_THREAD_KEYRING)
    if err != nil {
        return err
    }
    defer unlinkKey(kid)
    keyid := fmt.Sprintf(":%v:logon:%v", len(key), keyname)

    name := "crypttarget"
    uuid := "2f144136-b0de-4b51-b2eb-bd869cc39a6e"
    c := devmapper.CryptTable{
        Length:        60000 * 512, // size of the device in bytes
        Encryption:    "aes-xts-plain64",
        KeyID:         keyid,
        BackendDevice: "/dev/loop0",
        Flags:         []string{devmapper.CryptFlagAllowDiscards},
    }
    if err := devmapper.CreateAndLoad(name, uuid, c); err != nil {
        // handle error
    }
    defer devmapper.Remove(name)
}

func unlinkKey(kid int) {
	if _, err := unix.KeyctlInt(unix.KEYCTL_REVOKE, kid, 0, 0, 0); err != nil {
		fmt.Printf("key revoke: %v\n", err)
	}

	if _, err := unix.KeyctlInt(unix.KEYCTL_UNLINK, kid, unix.KEY_SPEC_THREAD_KEYRING, 0, 0); err != nil {
		fmt.Printf("key unlink, thread: %v\n", err)
	}

	// We added key to thread keyring only. But let's try to unlink the key from other keyrings as well just to be safe
	_, _ = unix.KeyctlInt(unix.KEYCTL_UNLINK, kid, unix.KEY_SPEC_PROCESS_KEYRING, 0, 0)
	_, _ = unix.KeyctlInt(unix.KEYCTL_UNLINK, kid, unix.KEY_SPEC_USER_KEYRING, 0, 0)
}
```

## Snapshot Target

The library supports device mapper snapshot targets for creating point-in-time copies of block devices.

### Snapshot-Origin Target

A snapshot-origin target represents the original device being protected by snapshots. When writes occur to the origin device, the original data is copied to a Copy-on-Write (CoW) device before being overwritten.

```go
func main() {
    name := "original-device"
    uuid := "550e8400-e29b-41d4-a716-446655440000"

    snapOrigin := devmapper.SnapshotOriginTable{
		Device: "/dev/sda",
		Start:  0,
		Length: 100 * 512, // Length of the block device, 100 sectors = 50KB
	}

    if err := devmapper.CreateAndLoad(name, uuid, 0, snapOrigin); err != nil {
        // handle error
    }
    defer devmapper.Remove(name)
}
```

### Snapshot Target

A snapshot target provides a consistent view of the origin device at the time the snapshot was created. Multiple independent snapshots can be created from the same origin.

```go
func main() {
    // Assume origin device /dev/mapper/origin-device exists
    // and CoW device /dev/mapper/cow-storage is allocated

    name := "snapshot-1"
    uuid := "550e8400-e29b-41d4-a716-446655440001"

    snapshot := devmapper.SnapshotTable{
        Length:         100 * 512,              // Same size as origin
        OriginDevice:   "/dev/mapper/origin-device",
        CowDevice:      "/dev/mapper/cow-storage",
        ChunkSize:      16,                     // Chunk size in sectors (must be power of 2)
        Persistent:     true,
    }

    if err := devmapper.CreateAndLoad(name, uuid, 0, snapshot); err != nil {
        // handle error
    }
    defer devmapper.Remove(name)

    // The snapshot is now active and ready to use
    // Reading from /dev/mapper/snapshot-1 returns data from the origin device
    // Writes to the origin device automatically trigger copy-on-write operations
}
```

## License

See [LICENSE](LICENSE).
