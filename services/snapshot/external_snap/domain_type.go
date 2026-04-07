package external

import (
	"fmt"

	"github.com/easy-cloud-Knet/KWS_Core/internal/snapshotlibvirt"
	"libvirt.org/go/libvirt"
)

type SnapshotDomain interface {
	CreateExternalSnapshot(snapshotXML string, opts externalSnapshotCreateExecOptions) (SnapshotHandle, error)
	ListAllSnapshots() ([]SnapshotHandle, error)
	StartBlockCommit(disk, base, top string) error
	BlockJobInfo(disk string) (externalBlockJobInfo, error)
	AbortBlockJobPivot(disk string) error
	UpdateDeviceConfig(deviceXML string) error
}

type externalBlockJobInfo struct {
	Cur uint64
	End uint64
}

type externalSnapshotCreateExecOptions struct {
	Live    bool
	Quiesce bool
	Atomic  bool
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

func (d *libvirtExternalSnapshotDomain) CreateExternalSnapshot(snapshotXML string, opts externalSnapshotCreateExecOptions) (SnapshotHandle, error) {
	if d == nil || d.domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	snapshot, err := snapshotlibvirt.CreateExternalSnapshot(d.domain, snapshotXML, snapshotlibvirt.ExternalSnapshotCreateOptions{
		Live:    opts.Live,
		Quiesce: opts.Quiesce,
		Atomic:  opts.Atomic,
	})
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

func (d *libvirtExternalSnapshotDomain) StartBlockCommit(disk, base, top string) error {
	if d == nil || d.domain == nil {
		return fmt.Errorf("nil domain")
	}

	return snapshotlibvirt.StartBlockCommit(d.domain, disk, base, top)
}

func (d *libvirtExternalSnapshotDomain) BlockJobInfo(disk string) (externalBlockJobInfo, error) {
	if d == nil || d.domain == nil {
		return externalBlockJobInfo{}, fmt.Errorf("nil domain")
	}

	job, err := d.domain.GetBlockJobInfo(disk, 0)
	if err != nil {
		return externalBlockJobInfo{}, err
	}

	return externalBlockJobInfo{Cur: job.Cur, End: job.End}, nil
}

func (d *libvirtExternalSnapshotDomain) AbortBlockJobPivot(disk string) error {
	if d == nil || d.domain == nil {
		return fmt.Errorf("nil domain")
	}

	return snapshotlibvirt.AbortBlockJobPivot(d.domain, disk)
}

func (d *libvirtExternalSnapshotDomain) UpdateDeviceConfig(deviceXML string) error {
	if d == nil || d.domain == nil {
		return fmt.Errorf("nil domain")
	}

	return snapshotlibvirt.UpdateDeviceConfig(d.domain, deviceXML)
}
