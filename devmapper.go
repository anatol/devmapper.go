package devmapper

import (
	"fmt"

	"golang.org/x/sys/unix"
)

type Table interface {
	startSector() uint64
	length() uint64
	targetType() string
	buildSpec() string // see https://wiki.gentoo.org/wiki/Device-mapper for examples of specs
}

var errNotImplemented = fmt.Errorf("Message functionality is not implemented")

func Create(name string, uuid string) error {
	return ioctl(unix.DM_DEV_CREATE, name, uuid, 0, false, nil)
}

func CreateAndLoad(name string, uuid string, tables ...Table) error {
	if err := Create(name, uuid); err != nil {
		return err
	}
	if err := Load(name, tables...); err != nil {
		return err
	}
	return Resume(name)
}

func Message(name string, sector int, message string) error {
	return errNotImplemented
}

func Suspend(name string) error {
	return ioctl(unix.DM_DEV_SUSPEND, name, "", unix.DM_SUSPEND_FLAG, false, nil)
}

func Resume(name string) error {
	return ioctl(unix.DM_DEV_SUSPEND, name, "", 0, true, nil)
}

func Load(name string, tables ...Table) error {
	return ioctl(unix.DM_TABLE_LOAD, name, "", 0, false, tables)
}

func Rename(old, new string) error {
	// primaryUdevEvent == true
	return errNotImplemented
}

func SetUUID(name, uuid string) error {
	return errNotImplemented
}

func Remove(name string) error {
	return ioctl(unix.DM_DEV_REMOVE, name, "", 0, true, nil)
}
