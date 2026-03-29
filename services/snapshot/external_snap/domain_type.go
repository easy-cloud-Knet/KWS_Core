package external

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

type externalSnapshotDomain interface {
	IsActive() (bool, error)
	CreateExternalSnapshot(snapshotXML string, opts externalSnapshotCreateExecOptions) (externalSnapshotHandle, error)
	ListAllSnapshots() ([]externalSnapshotHandle, error)
	StartBlockCommit(disk, base, top string) error
	BlockJobInfo(disk string) (externalBlockJobInfo, error)
	AbortBlockJobPivot(disk string) error
	UpdateDeviceConfig(deviceXML string) error
	UUIDString() (string, error)
	XMLDesc() (string, error)
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

type externalSnapshotHandle interface {
	Name() (string, error)
	XMLDesc() (string, error)
	Delete() error
	Free() error
}

type libvirtExternalSnapshotDomain struct {
	domain *libvirt.Domain
}

func newExternalSnapshotDomain(domain *libvirt.Domain) externalSnapshotDomain {
	return &libvirtExternalSnapshotDomain{domain: domain}
}

func (d *libvirtExternalSnapshotDomain) IsActive() (bool, error) {
	if d == nil || d.domain == nil {
		return false, fmt.Errorf("nil domain")
	}
	return d.domain.IsActive()
}

func (d *libvirtExternalSnapshotDomain) CreateExternalSnapshot(snapshotXML string, opts externalSnapshotCreateExecOptions) (externalSnapshotHandle, error) {
	if d == nil || d.domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	flags := libvirt.DOMAIN_SNAPSHOT_CREATE_DISK_ONLY
	if opts.Live {
		flags |= libvirt.DOMAIN_SNAPSHOT_CREATE_LIVE
	}
	if opts.Quiesce {
		flags |= libvirt.DOMAIN_SNAPSHOT_CREATE_QUIESCE
	}
	if opts.Atomic {
		flags |= libvirt.DOMAIN_SNAPSHOT_CREATE_ATOMIC
	}

	snapshot, err := d.domain.CreateSnapshotXML(snapshotXML, flags)
	if err != nil {
		return nil, err
	}

	return &libvirtExternalSnapshotHandle{snapshot: snapshot}, nil
}

func (d *libvirtExternalSnapshotDomain) UUIDString() (string, error) {
	if d == nil || d.domain == nil {
		return "", fmt.Errorf("nil domain")
	}
	return d.domain.GetUUIDString()
}

func (d *libvirtExternalSnapshotDomain) ListAllSnapshots() ([]externalSnapshotHandle, error) {
	if d == nil || d.domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	snaps, err := d.domain.ListAllSnapshots(0)
	if err != nil {
		return nil, err
	}

	handles := make([]externalSnapshotHandle, 0, len(snaps))
	for i := range snaps {
		handles = append(handles, &libvirtExternalSnapshotHandle{snapshot: &snaps[i]})
	}

	return handles, nil
}

func (d *libvirtExternalSnapshotDomain) StartBlockCommit(disk, base, top string) error {
	if d == nil || d.domain == nil {
		return fmt.Errorf("nil domain")
	}

	flags := libvirt.DOMAIN_BLOCK_COMMIT_ACTIVE | libvirt.DOMAIN_BLOCK_COMMIT_DELETE
	return d.domain.BlockCommit(disk, base, top, 0, flags)
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

	return d.domain.BlockJobAbort(disk, libvirt.DOMAIN_BLOCK_JOB_ABORT_PIVOT)
}

func (d *libvirtExternalSnapshotDomain) UpdateDeviceConfig(deviceXML string) error {
	if d == nil || d.domain == nil {
		return fmt.Errorf("nil domain")
	}

	return d.domain.UpdateDeviceFlags(deviceXML, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG)
}

func (d *libvirtExternalSnapshotDomain) XMLDesc() (string, error) {
	if d == nil || d.domain == nil {
		return "", fmt.Errorf("nil domain")
	}
	return d.domain.GetXMLDesc(0)
}

type libvirtExternalSnapshotHandle struct {
	snapshot *libvirt.DomainSnapshot
}

func (s *libvirtExternalSnapshotHandle) Name() (string, error) {
	if s == nil || s.snapshot == nil {
		return "", fmt.Errorf("nil snapshot")
	}
	return s.snapshot.GetName()
}

func (s *libvirtExternalSnapshotHandle) XMLDesc() (string, error) {
	if s == nil || s.snapshot == nil {
		return "", fmt.Errorf("nil snapshot")
	}
	return s.snapshot.GetXMLDesc(0)
}

func (s *libvirtExternalSnapshotHandle) Delete() error {
	if s == nil || s.snapshot == nil {
		return fmt.Errorf("nil snapshot")
	}
	return s.snapshot.Delete(0)
}

func (s *libvirtExternalSnapshotHandle) Free() error {
	if s == nil || s.snapshot == nil {
		return nil
	}
	return s.snapshot.Free()
}
