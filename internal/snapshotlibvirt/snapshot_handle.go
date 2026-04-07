package snapshotlibvirt

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

type LibvirtSnapshotHandle struct {
	snapshot *libvirt.DomainSnapshot
}

func NewSnapshotHandle(snapshot *libvirt.DomainSnapshot) *LibvirtSnapshotHandle {
	return &LibvirtSnapshotHandle{snapshot: snapshot}
}

func (s *LibvirtSnapshotHandle) Name() (string, error) {
	if s == nil || s.snapshot == nil {
		return "", fmt.Errorf("nil snapshot")
	}
	return s.snapshot.GetName()
}

func (s *LibvirtSnapshotHandle) Delete() error {
	if s == nil || s.snapshot == nil {
		return fmt.Errorf("nil snapshot")
	}
	return s.snapshot.Delete(0)
}

func (s *LibvirtSnapshotHandle) Revert() error {
	if s == nil || s.snapshot == nil {
		return fmt.Errorf("nil snapshot")
	}
	return s.snapshot.RevertToSnapshot(0)
}

func (s *LibvirtSnapshotHandle) XMLDesc() (string, error) {
	if s == nil || s.snapshot == nil {
		return "", fmt.Errorf("nil snapshot")
	}
	return s.snapshot.GetXMLDesc(0)
}

func (s *LibvirtSnapshotHandle) Free() error {
	if s == nil || s.snapshot == nil {
		return nil
	}
	return s.snapshot.Free()
}
