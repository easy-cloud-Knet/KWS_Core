package snapshotlibvirt

import "libvirt.org/go/libvirt"

type ExternalSnapshotCreateOptions struct {
	Live    bool
	Quiesce bool
	Atomic  bool
}

func CreateExternalSnapshot(domain *libvirt.Domain, snapshotXML string, opts ExternalSnapshotCreateOptions) (*libvirt.DomainSnapshot, error) {
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

	return domain.CreateSnapshotXML(snapshotXML, flags)
}

func StartBlockCommit(domain *libvirt.Domain, disk, base, top string) error {
	flags := libvirt.DOMAIN_BLOCK_COMMIT_ACTIVE | libvirt.DOMAIN_BLOCK_COMMIT_DELETE
	return domain.BlockCommit(disk, base, top, 0, flags)
}

func AbortBlockJobPivot(domain *libvirt.Domain, disk string) error {
	return domain.BlockJobAbort(disk, libvirt.DOMAIN_BLOCK_JOB_ABORT_PIVOT)
}

func UpdateDeviceConfig(domain *libvirt.Domain, deviceXML string) error {
	return domain.UpdateDeviceFlags(deviceXML, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG)
}
