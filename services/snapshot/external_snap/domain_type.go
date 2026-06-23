package external

import (
	"fmt"

	"github.com/easy-cloud-Knet/KWS_Core/internal/snapshotlibvirt"
	"libvirt.org/go/libvirt"
)

type SnapshotDomain interface {
	RegisterExternalSnapshot(snapshotXML string) (SnapshotHandle, error)
	ListAllSnapshots() ([]SnapshotHandle, error)
	UpdateDeviceConfig(deviceXML string) error
}

type SnapshotHandle interface {
	Name() (string, error)
	XMLDesc() (string, error)
	Delete() error
	Free() error
}

type libvirtExternalSnapshotDomain struct {
	domain *libvirt.Domain
}

func newExternalSnapshotDomain(domain *libvirt.Domain) SnapshotDomain {
	return &libvirtExternalSnapshotDomain{domain: domain}
}

func (d *libvirtExternalSnapshotDomain) RegisterExternalSnapshot(snapshotXML string) (SnapshotHandle, error) {
	if d == nil || d.domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	snapshot, err := snapshotlibvirt.RegisterExternalSnapshot(d.domain, snapshotXML)
	if err != nil {
		return nil, err
	}

	return snapshotlibvirt.NewSnapshotHandle(snapshot), nil
}

func (d *libvirtExternalSnapshotDomain) ListAllSnapshots() ([]SnapshotHandle, error) {
	if d == nil || d.domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	snaps, err := d.domain.ListAllSnapshots(0)
	if err != nil {
		return nil, err
	}

	handles := make([]SnapshotHandle, 0, len(snaps))
	for i := range snaps {
		handles = append(handles, snapshotlibvirt.NewSnapshotHandle(&snaps[i]))
	}

	return handles, nil
}

func (d *libvirtExternalSnapshotDomain) UpdateDeviceConfig(deviceXML string) error {
	if d == nil || d.domain == nil {
		return fmt.Errorf("nil domain")
	}

	return snapshotlibvirt.UpdateDeviceConfig(d.domain, deviceXML)
}
