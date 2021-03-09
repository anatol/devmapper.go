package devmapper

import (
	"fmt"
	"golang.org/x/sys/unix"
	"unsafe"
)

// device size calculation happens in sector units of size 512
const SectorSize = 512

type Table interface {
	startSector() uint64
	length() uint64
	targetType() string
	buildSpec() string // see https://wiki.gentoo.org/wiki/Device-mapper for examples of specs
}

var errNotImplemented = fmt.Errorf("Message functionality is not implemented")

// Create creates a new device. No table will be loaded. The device will be in
// suspended state. Any IO to this device will fail.
func Create(name string, uuid string) error {
	return ioctlTable(unix.DM_DEV_CREATE, name, uuid, 0, false, nil)
}

// CreateAndLoad creates, loads the provided tables and resumes the device.
func CreateAndLoad(name string, uuid string, tables ...Table) error {
	if err := Create(name, uuid); err != nil {
		return err
	}
	if err := Load(name, tables...); err != nil {
		return err
	}
	return Resume(name)
}

// Message passes a message string to the target at specific offset of a device.
func Message(name string, sector int, message string) error {
	return errNotImplemented
}

// Suspend suspends the given device.
func Suspend(name string) error {
	return ioctlTable(unix.DM_DEV_SUSPEND, name, "", unix.DM_SUSPEND_FLAG, false, nil)
}

// Resume resumes the given device.
func Resume(name string) error {
	return ioctlTable(unix.DM_DEV_SUSPEND, name, "", 0, true, nil)
}

// Load loads given table into the device
func Load(name string, tables ...Table) error {
	return ioctlTable(unix.DM_TABLE_LOAD, name, "", 0, false, tables)
}

// Rename renames the device
func Rename(old, new string) error {
	// primaryUdevEvent == true
	return errNotImplemented
}

// SetUUID sets uuid for a given device
func SetUUID(name, uuid string) error {
	return errNotImplemented
}

// Remove remove the device and destroys its tables.
func Remove(name string) error {
	return ioctlTable(unix.DM_DEV_REMOVE, name, "", 0, true, nil)
}

// GetVersion returns version for the dm-mapper kernel interface
func GetVersion() (major, minor, patch uint32, err error) {
	data := make([]byte, unix.SizeofDmIoctl)
	ioctlData := (*unix.DmIoctl)(unsafe.Pointer(&data[0]))
	ioctlData.Version = [...]uint32{4, 0, 0} // minimum required version
	ioctlData.Data_size = unix.SizeofDmIoctl
	ioctlData.Data_start = unix.SizeofDmIoctl

	if err := ioctl(unix.DM_VERSION, data); err != nil {
		return 0, 0, 0, err
	}

	return ioctlData.Version[0], ioctlData.Version[1], ioctlData.Version[2], nil
}
