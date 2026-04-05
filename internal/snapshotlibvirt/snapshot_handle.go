package snapshotlibvirt

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

type SnapshotHandle interface {
	Name() (string, error)
	Delete() error
	Revert() error
	Free() error
}

type libvirtSnapshotHandle struct {
	snapshot *libvirt.DomainSnapshot
}

func NewSnapshotHandle(snapshot *libvirt.DomainSnapshot) SnapshotHandle {
	return &libvirtSnapshotHandle{snapshot: snapshot}
}

func (s *libvirtSnapshotHandle) Name() (string, error) {
	if s == nil || s.snapshot == nil {
		return "", fmt.Errorf("nil snapshot")
	}
	return s.snapshot.GetName()
}

func (s *libvirtSnapshotHandle) Delete() error {
	if s == nil || s.snapshot == nil {
		return fmt.Errorf("nil snapshot")
	}
	return s.snapshot.Delete(0)
}

func (s *libvirtSnapshotHandle) Revert() error {
	if s == nil || s.snapshot == nil {
		return fmt.Errorf("nil snapshot")
	}
	return s.snapshot.RevertToSnapshot(0)
}

func (s *libvirtSnapshotHandle) Free() error {
	if s == nil || s.snapshot == nil {
		return nil
	}
	return s.snapshot.Free()
}
