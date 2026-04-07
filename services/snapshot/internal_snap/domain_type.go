package internal

import (
	"fmt"

	"github.com/easy-cloud-Knet/KWS_Core/internal/snapshotlibvirt"
	"libvirt.org/go/libvirt"
)

type snapshotDomain interface {
	CreateSnapshot(snapshotXML string, opts snapshotCreateOptions) (snapshotHandle, error)
	ListAllSnapshots() ([]snapshotHandle, error)
}

type snapshotCreateOptions struct {
	Quiesce bool
}

type snapshotHandle interface {
	Name() (string, error)
	Delete() error
	Revert() error
	Free() error
}

type libvirtSnapshotDomain struct {
	domain *libvirt.Domain
}

func newInternalSnapshotDomain(domain *libvirt.Domain) snapshotDomain {
	return &libvirtSnapshotDomain{domain: domain}
}

func (d *libvirtSnapshotDomain) CreateSnapshot(snapshotXML string, opts snapshotCreateOptions) (snapshotHandle, error) {
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

	return snapshotlibvirt.NewSnapshotHandle(snap), nil
}

func (d *libvirtSnapshotDomain) ListAllSnapshots() ([]snapshotHandle, error) {
	if d == nil || d.domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	snaps, err := d.domain.ListAllSnapshots(0)
	if err != nil {
		return nil, err
	}

	handles := make([]snapshotHandle, 0, len(snaps))
	for i := range snaps {
		handles = append(handles, snapshotlibvirt.NewSnapshotHandle(&snaps[i]))
	}

	return handles, nil
}
