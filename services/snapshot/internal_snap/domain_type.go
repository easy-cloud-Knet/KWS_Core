package internal

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

type internalSnapshotDomain interface {
	CreateSnapshot(snapshotXML string, opts internalSnapshotCreateExecOptions) (internalSnapshotHandle, error)
	ListAllSnapshots() ([]internalSnapshotHandle, error)
}

type internalSnapshotCreateExecOptions struct {
	Quiesce bool
}

type internalSnapshotHandle interface {
	Name() (string, error)
	Delete() error
	Revert() error
	Free() error
}

type libvirtInternalSnapshotDomain struct {
	domain *libvirt.Domain
}

func newInternalSnapshotDomain(domain *libvirt.Domain) internalSnapshotDomain {
	return &libvirtInternalSnapshotDomain{domain: domain}
}

func (d *libvirtInternalSnapshotDomain) CreateSnapshot(snapshotXML string, opts internalSnapshotCreateExecOptions) (internalSnapshotHandle, error) {
	if d == nil || d.domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	flags := libvirt.DomainSnapshotCreateFlags(0)
	if opts.Quiesce {
		flags |= libvirt.DOMAIN_SNAPSHOT_CREATE_QUIESCE
	}

	snap, err := d.domain.CreateSnapshotXML(snapshotXML, flags)
	if err != nil {
		return nil, err
	}

	return &libvirtInternalSnapshotHandle{snapshot: snap}, nil
}

func (d *libvirtInternalSnapshotDomain) ListAllSnapshots() ([]internalSnapshotHandle, error) {
	if d == nil || d.domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	snaps, err := d.domain.ListAllSnapshots(0)
	if err != nil {
		return nil, err
	}

	handles := make([]internalSnapshotHandle, 0, len(snaps))
	for i := range snaps {
		handles = append(handles, &libvirtInternalSnapshotHandle{snapshot: &snaps[i]})
	}

	return handles, nil
}

type libvirtInternalSnapshotHandle struct {
	snapshot *libvirt.DomainSnapshot
}

func (s *libvirtInternalSnapshotHandle) Name() (string, error) {
	if s == nil || s.snapshot == nil {
		return "", fmt.Errorf("nil snapshot")
	}
	return s.snapshot.GetName()
}

func (s *libvirtInternalSnapshotHandle) Delete() error {
	if s == nil || s.snapshot == nil {
		return fmt.Errorf("nil snapshot")
	}
	return s.snapshot.Delete(0)
}

func (s *libvirtInternalSnapshotHandle) Revert() error {
	if s == nil || s.snapshot == nil {
		return fmt.Errorf("nil snapshot")
	}
	return s.snapshot.RevertToSnapshot(0)
}

func (s *libvirtInternalSnapshotHandle) Free() error {
	if s == nil || s.snapshot == nil {
		return nil
	}
	return s.snapshot.Free()
}
