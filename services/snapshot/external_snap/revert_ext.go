package external

import (
	"fmt"
	"os"
	"path/filepath"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"github.com/easy-cloud-Knet/KWS_Core/internal/qemuimg"
)

func RevertExternalSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	active, err := domain.Domain.IsActive()
	if err == nil && active {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("external snapshot revert requires the domain to be shut down"))
	}

	xmlDesc, err := domain.Domain.GetXMLDesc(0)
	if err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to get domain xml: %w", err))
	}

	return revertExternalSnapshot(newExternalSnapshotDomain(domain.Domain), qemuimg.New(), xmlDesc, snapName)
}

func revertExternalSnapshot(domain SnapshotDomain, qimg qemuimg.Runner, domainXMLDesc, snapName string) error {
	if domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}
	if snapName == "" {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot name required"))
	}

	snaps, err := domain.ListAllSnapshots()
	if err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list snapshots: %w", err))
	}
	defer freeSnapshotHandles(snaps)

	target, err := findExternalSnapshotByName(snaps, snapName)
	if err != nil {
		return err
	}
	if target == nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("snapshot %s not found", snapName))
	}

	// snapOverlays maps diskName → overlay file path recorded in the snapshot
	snapOverlays, err := extractExternalSnapshotSources(target)
	if err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to extract snapshot sources: %w", err))
	}
	if len(snapOverlays) == 0 {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("snapshot %s has no external disk sources", snapName))
	}

	disks, err := listFileDisksFromXMLDesc(domainXMLDesc)
	if err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list file disks: %w", err))
	}

	updated := false
	for _, d := range disks {
		snapOverlay, ok := snapOverlays[d.TargetDev]
		if !ok || snapOverlay == "" {
			continue
		}

		// Resolve the backing file of the snapshot overlay — this is the state at snapshot time.
		backingFile, backingFormat, err := qimg.Info(snapOverlay)
		if err != nil {
			return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to query backing file for disk %s: %w", d.TargetDev, err))
		}
		if backingFile == "" {
			return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("snapshot overlay for disk %s has no backing file", d.TargetDev))
		}
		if backingFormat == "" {
			backingFormat = "qcow2"
		}

		// Create a fresh writable overlay backed by the snapshot's backing file.
		// Layout: <root>/<uuid>/working/<disk>.qcow2
		workingPath, err := workingDiskPath(snapOverlay, d.TargetDev)
		if err != nil {
			return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("invalid working disk path for disk %s: %w", d.TargetDev, err))
		}
		if err := os.MkdirAll(filepath.Dir(workingPath), 0755); err != nil {
			return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to create working directory for disk %s: %w", d.TargetDev, err))
		}

		if err := qimg.ReplaceOverlay(backingFile, backingFormat, workingPath); err != nil {
			return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to replace working overlay for disk %s: %w", d.TargetDev, err))
		}

		diskXML := buildDiskDeviceXML(d, workingPath)
		if err := domain.UpdateDeviceConfig(diskXML); err != nil {
			return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to update disk %s: %w", d.TargetDev, err))
		}

		updated = true
	}

	if !updated {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("no matching disks found to revert"))
	}

	return nil
}
