package external

import (
	"fmt"
	"time"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func MergeExternalSnapshot(domain *domCon.Domain, targetDisk string) ([]string, error) {
	if domain == nil || domain.Domain == nil {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	return mergeExternalSnapshot(newExternalSnapshotDomain(domain.Domain), targetDisk)
}

func mergeExternalSnapshot(domain externalSnapshotDomain, targetDisk string) ([]string, error) {
	if domain == nil {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	active, err := domain.IsActive()
	if err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to check domain state: %w", err))
	}
	if !active {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("external snapshot merge requires the domain to be running"))
	}

	disks, err := listFileDisks(domain)
	if err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list file disks: %w", err))
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

		if err := domain.StartBlockCommit(d.TargetDev, backingSource, d.Source); err != nil {
			return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to start block commit on disk %s: %w", d.TargetDev, err))
		}

		if err := waitBlockJobReady(domain, d.TargetDev, 2*time.Minute); err != nil {
			return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("block commit did not complete for disk %s: %w", d.TargetDev, err))
		}

		if err := domain.AbortBlockJobPivot(d.TargetDev); err != nil {
			return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to pivot disk %s after commit: %w", d.TargetDev, err))
		}

		diskXML := buildDiskDeviceXML(d, backingSource)
		if err := domain.UpdateDeviceConfig(diskXML); err != nil {
			return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to update disk %s after merge: %w", d.TargetDev, err))
		}

		merged = append(merged, d.TargetDev)
	}

	if targetDisk != "" && len(merged) == 0 {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("disk %s is not mergeable or has no external backing chain", targetDisk))
	}
	if len(merged) == 0 {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("no mergeable external snapshot disks found"))
	}

	return merged, nil
}
