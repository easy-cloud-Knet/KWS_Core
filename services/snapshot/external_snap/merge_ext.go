package external

import (
	"fmt"
	"time"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"libvirt.org/go/libvirt"
)

func MergeExternalSnapshot(domain *domCon.Domain, targetDisk string) ([]string, error) {
	if domain == nil || domain.Domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	active, err := domain.Domain.IsActive()
	if err != nil {
		return nil, fmt.Errorf("failed to check domain state: %w", err)
	}
	if !active {
		return nil, fmt.Errorf("external snapshot merge requires the domain to be running")
	}

	disks, err := listFileDisks(domain)
	if err != nil {
		return nil, err
	}

	merged := make([]string, 0, len(disks))
	for _, d := range disks {
		if targetDisk != "" && d.TargetDev != targetDisk {
			continue
		}

		backingSource := d.BackingSource

		if backingSource == "" {
			continue
		}

		flags := libvirt.DOMAIN_BLOCK_COMMIT_ACTIVE | libvirt.DOMAIN_BLOCK_COMMIT_DELETE
		if err := domain.Domain.BlockCommit(d.TargetDev, backingSource, d.Source, 0, flags); err != nil {
			return nil, fmt.Errorf("failed to start block commit on disk %s: %w", d.TargetDev, err)
		}

		if err := waitBlockJobReady(domain.Domain, d.TargetDev, 2*time.Minute); err != nil {
			return nil, err
		}

		if err := domain.Domain.BlockJobAbort(d.TargetDev, libvirt.DOMAIN_BLOCK_JOB_ABORT_PIVOT); err != nil {
			return nil, fmt.Errorf("failed to pivot disk %s after commit: %w", d.TargetDev, err)
		}

		diskXML := buildDiskDeviceXML(d, backingSource)
		if err := domain.Domain.UpdateDeviceFlags(diskXML, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG); err != nil {
			return nil, fmt.Errorf("failed to update disk %s after merge: %w", d.TargetDev, err)
		}

		merged = append(merged, d.TargetDev)
	}

	if targetDisk != "" && len(merged) == 0 {
		return nil, fmt.Errorf("disk %s is not mergeable or has no external backing chain", targetDisk)
	}
	if len(merged) == 0 {
		return nil, fmt.Errorf("no mergeable external snapshot disks found")
	}

	return merged, nil
}
