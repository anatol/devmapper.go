# devmapper.go

A pure-Go library for Linux device mapper target management.

## What is device mapper?

Device mapper is a Linux kernel framework that provides a way to create virtual block devices by mapping them onto one or more underlying block devices. It lives in the kernel and is controlled from userspace via the `/dev/mapper/control` ioctl interface. The `dmsetup` command-line tool is the most common way to interact with it, but `devmapper.go` lets you do the same directly from Go without shelling out.

Device mapper works by defining **tables** that describe how the virtual device's sectors map to underlying storage. Each table specifies a **target type** (e.g. `linear`, `crypt`, `striped`) that determines how I/O is transformed. Multiple tables can be stacked to build complex storage configurations — for example, encrypting data on top of a striped RAID set.

Device mapper is the foundation behind many Linux storage tools:

- **LVM2** — logical volume management
- **LUKS/cryptsetup** — full-disk encryption
- **dm-verity** — read-only integrity verification (used by Android and Chrome OS)
- **Docker/containerd** — thin provisioning for container images
- **multipath-tools** — SAN multipath I/O

## Installation

```
go get github.com/anatol/devmapper.go
```

## Supported targets

| Target | Go type | Description |
|---|---|---|
| `linear` | `LinearTable` | Maps a range of sectors directly to another block device |
| `striped` | `StripeTable` | Stripes I/O across multiple devices (RAID-0) |
| `mirror` | `MirrorTable` | Mirrors I/O across multiple devices (RAID-1) |
| `raid` | `RaidTable` | Various RAID levels using the md framework (raid0, raid1, raid4, raid5, raid6, raid10) |
| `crypt` | `CryptTable` | Encrypts data using the kernel crypto API (used by LUKS/cryptsetup) |
| `verity` | `VerityTable` | Read-only integrity verification using a hash tree |
| `integrity` | `IntegrityTable` | Block-level data integrity protection (dm-integrity) |
| `snapshot` | `SnapshotTable` | Copy-on-write snapshots of a block device |
| `snapshot-origin` | `SnapshotOriginTable` | Marks a device as the origin for snapshots |
| `snapshot-merge` | `SnapshotMergeTable` | Merges a snapshot back into its origin |
| `thin-pool` | `ThinPoolTable` | Manages a pool of blocks shared between thin devices |
| `thin` | `ThinTable` | Thin-provisioned device backed by a thin-pool |
| `cache` | `CacheTable` | Uses a fast device (e.g. SSD) to cache a slower origin device |
| `writecache` | `WritecacheTable` | Caches writes on a fast device, flushes later to slower storage |
| `era` | `EraTable` | Tracks which blocks were written within a given period |
| `delay` | `DelayTable` | Delays I/O for testing purposes |
| `flakey` | `FlakeyTable` | Simulates unreliable device (periodic I/O errors) for testing |
| `dust` | `DustTable` | Simulates failing sectors for testing |
| `log-writes` | `LogWritesTable` | Logs all write operations to a separate device for crash testing |
| `switch` | `SwitchTable` | Switches between multiple backend devices per region |
| `unstriped` | `UnstripedTable` | Extracts a single stripe from a striped device |
| `multipath` | `MultipathTable` | I/O failover and load balancing across multiple paths |
| `error` | `ErrorTable` | Always fails I/O with EIO (useful for testing) |
| `zero` | `ZeroTable` | Returns zeros on read, discards writes (like `/dev/zero`) |

## Examples

### Create an encrypted device (dm-crypt)

```go
name := "mycrypt"
uuid := "2f144136-b0de-4b51-b2eb-bd869cc39a6e"
key := make([]byte, 32)

c := devmapper.CryptTable{
    Length:        60000 * devmapper.SectorSize,
    Encryption:    "aes-xts-plain64",
    Key:           key,
    BackendDevice: "/dev/loop0",
    Flags:         []string{devmapper.CryptFlagAllowDiscards},
}
if err := devmapper.CreateAndLoad(name, uuid, 0, c); err != nil {
    log.Fatal(err)
}
defer devmapper.Remove(name)

// /dev/mapper/mycrypt is now available
```

### Create an encrypted device using Linux kernel keyring

```go
// Load key into kernel keyring
keyname := fmt.Sprintf("cryptsetup:%s-d%d", uuid, luksDigestId)
kid, err := unix.AddKey("logon", keyname, key, unix.KEY_SPEC_THREAD_KEYRING)
if err != nil {
    log.Fatal(err)
}
keyid := fmt.Sprintf(":%v:logon:%v", len(key), keyname)

c := devmapper.CryptTable{
    Length:        60000 * devmapper.SectorSize,
    Encryption:    "aes-xts-plain64",
    KeyID:         keyid,
    BackendDevice: "/dev/loop0",
    Flags:         []string{devmapper.CryptFlagAllowDiscards},
}
if err := devmapper.CreateAndLoad(name, uuid, 0, c); err != nil {
    log.Fatal(err)
}
defer devmapper.Remove(name)
```

### Create a linear mapping

```go
l := devmapper.LinearTable{
    Length:        15 * devmapper.SectorSize,
    BackendDevice: "/dev/loop0",
}
if err := devmapper.CreateAndLoad("mylinear", "some-uuid", 0, l); err != nil {
    log.Fatal(err)
}
defer devmapper.Remove("mylinear")
```

### Create a striped (RAID-0) device

```go
s := devmapper.StripeTable{
    Length:    totalSize,
    ChunkSize: 64 * 1024, // 64 KiB chunks
    Devices: []devmapper.StripeDevice{
        {Device: "/dev/loop0"},
        {Device: "/dev/loop1"},
        {Device: "/dev/loop2"},
    },
}
if err := devmapper.CreateAndLoad("mystripe", "some-uuid", 0, s); err != nil {
    log.Fatal(err)
}
defer devmapper.Remove("mystripe")
```

### Create a read-only verity device

```go
v := devmapper.VerityTable{
    Length:        dataSize,
    HashType:      1,
    DataDevice:    "/dev/loop0",
    HashDevice:    "/dev/loop1",
    DataBlockSize: 4096,
    HashBlockSize: 4096,
    NumDataBlocks: 512,
    Algorithm:     "sha256",
    Salt:          salt,
    Digest:        rootHash,
}
if err := devmapper.CreateAndLoad("myverity", "some-uuid", devmapper.ReadOnlyFlag, v); err != nil {
    log.Fatal(err)
}
defer devmapper.Remove("myverity")
```

### Query device mapper state

```go
// List all device mapper devices
items, err := devmapper.List()

// Get device info
info, err := devmapper.InfoByName("mydevice")
fmt.Printf("Name: %s, UUID: %s, Open count: %d\n", info.Name, info.UUID, info.OpenCount)

// Get kernel device mapper version
major, minor, patch, err := devmapper.GetVersion()
```

## Userspace volumes

In addition to creating kernel device mapper targets, `devmapper.go` supports **userspace volumes** — a way to read and write data using the same table definitions but without involving the kernel device mapper framework. Instead, all data processing (e.g. encryption, striping) happens in Go userspace.

This is useful for:

- Reading/writing encrypted or striped data on systems where you cannot load device mapper kernel modules (e.g. containers, CI)
- Unit testing without root privileges or device mapper setup
- Offline processing of disk images

```go
c := devmapper.CryptTable{
    Length:        size,
    Encryption:    "aes-xts-plain64",
    Key:           key,
    BackendDevice: "/path/to/encrypted.img",
}

vol, err := devmapper.OpenUserspaceVolume(os.O_RDWR, 0, c)
if err != nil {
    log.Fatal(err)
}
defer vol.Close()

// Read decrypted data
buf := make([]byte, 4096)
_, err = vol.ReadAt(buf, 0)

// Write data (gets encrypted transparently)
_, err = vol.WriteAt(plaintext, 0)
```

Userspace volumes implement the `io.ReaderAt`, `io.WriterAt`, and `io.Closer` interfaces. Reads and writes must be aligned to `devmapper.SectorSize` (512 bytes). Multiple tables can be passed to `OpenUserspaceVolume` to handle devices with multiple table regions.

The following targets have userspace volume support: `linear`, `striped`, `mirror`, `crypt`, `zero`, `error`, `delay`, `flakey`, `dust`, `log-writes`, `unstriped`, `multipath`.

## License

See [LICENSE](LICENSE).
