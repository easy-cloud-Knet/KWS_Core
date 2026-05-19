package external

import (
	"fmt"
	"os"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func MergeExternalSnapshot(domain *domCon.Domain, targetDisk string) ([]string, error) {
	if domain == nil || domain.Domain == nil {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	active, err := domain.Domain.IsActive()
	if err == nil && active {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("offline merge requires the domain to be shut down"))
	}

	xmlDesc, err := domain.Domain.GetXMLDesc(0)
	if err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to get domain xml: %w", err))
	}

	return mergeExternalSnapshot(newExternalSnapshotDomain(domain.Domain), newQemuImg(), xmlDesc, targetDisk)
}

func mergeExternalSnapshot(domain SnapshotDomain, qimg QemuImg, domainXMLDesc, targetDisk string) ([]string, error) {
	if domain == nil {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	disks, err := listFileDisksFromXMLDesc(domainXMLDesc)
	if err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list file disks: %w", err))
	}

	snaps, err := domain.ListAllSnapshots()
	if err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list snapshots: %w", err))
	}
	defer freeSnapshotHandles(snaps)

	type mergeTarget struct {
		disk       diskInfo
		originPath string
		overlays   []string
	}
	var targets []mergeTarget

	// Phase 1: commit each overlay chain into the VM's origin disk.
	// qemu-img commit -b <origin> <top_overlay> collapses all layers above origin
	// into origin. Only origin needs a write lock; the base image is read-only.
	for _, d := range disks {
		if targetDisk != "" && d.TargetDev != targetDisk {
			continue
		}

		backingFile, _, err := qimg.Info(d.Source)
		if err != nil {
			return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to query backing file for disk %s: %w", d.TargetDev, err))
		}
		if backingFile == "" {
			// Disk has no backing chain — nothing to merge.
			continue
		}

		originPath, overlays, err := findOriginAndOverlays(qimg, d.Source)
		if err != nil {
			return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to resolve disk chain for %s: %w", d.TargetDev, err))
		}

		if err := qimg.Commit(d.Source, originPath); err != nil {
			return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to commit disk %s into origin: %w", d.TargetDev, err))
		}

		targets = append(targets, mergeTarget{d, originPath, overlays})
	}

	if targetDisk != "" && len(targets) == 0 {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("disk %s has no backing store to merge", targetDisk))
	}
	if len(targets) == 0 {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("no mergeable disks found"))
	}

	// Phase 2: delete snapshot metadata while the domain still points to the
	// overlay chain so libvirt's reachability validation passes.
	deleteSnapshotMetadataLeafFirst(snaps)

	// Phase 3: update domain config to origin and remove the overlay files.
	merged := make([]string, 0, len(targets))
	for _, t := range targets {
		originDisk := diskInfo{
			TargetDev:  t.disk.TargetDev,
			TargetBus:  t.disk.TargetBus,
			Driver:     t.disk.Driver,
			DriverName: t.disk.DriverName,
		}
		diskXML := buildDiskDeviceXML(originDisk, t.originPath)
		if err := domain.UpdateDeviceConfig(diskXML); err != nil {
			return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to update disk %s after merge: %w", t.disk.TargetDev, err))
		}

		for _, overlayPath := range t.overlays {
			if err := os.Remove(overlayPath); err != nil && !os.IsNotExist(err) {
				return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to remove overlay %s: %w", overlayPath, err))
			}
		}

		merged = append(merged, t.disk.TargetDev)
	}

	return merged, nil
}

// deleteSnapshotMetadataLeafFirst deletes libvirt snapshot records leaf-first so
// that parent records become deletable once all their children are gone.
// Errors are ignored per iteration; the loop stops when no further progress is made.
func deleteSnapshotMetadataLeafFirst(snaps []SnapshotHandle) {
	deleted := make([]bool, len(snaps))
	for pass := 0; pass <= len(snaps); pass++ {
		progress := false
		for i, s := range snaps {
			if deleted[i] {
				continue
			}
			if s.Delete() == nil {
				deleted[i] = true
				progress = true
			}
		}
		if !progress {
			break
		}
	}
}
