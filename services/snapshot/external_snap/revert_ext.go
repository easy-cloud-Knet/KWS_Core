package external

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"libvirt.org/go/libvirt"
)

func RevertExternalSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}
	if snapName == "" {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot name required"))
	}

	active, err := domain.Domain.IsActive()
	if err == nil && active {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("external snapshot revert requires the domain to be shut down"))
	}

	snaps, err := domain.Domain.ListAllSnapshots(0)
	if err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list snapshots: %w", err))
	}
	defer func() {
		for _, s := range snaps {
			s.Free()
		}
	}()

	var target *libvirt.DomainSnapshot
	for i := range snaps {
		name, err := snaps[i].GetName()
		if err != nil || name != snapName {
			continue
		}
		isExternal, err := isExternalSnapshot(&snaps[i])
		if err != nil {
			return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to inspect snapshot %s: %w", snapName, err))
		}
		if !isExternal {
			return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot %s is not external", snapName))
		}
		target = &snaps[i]
		break
	}

	if target == nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("snapshot %s not found", snapName))
	}

	targetSources, err := extractExternalSnapshotSources(target)
	if err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to extract snapshot sources: %w", err))
	}
	if len(targetSources) == 0 {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("snapshot %s has no external disk sources", snapName))
	}

	disks, err := listFileDisks(domain)
	if err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list file disks: %w", err))
	}

	updated := false
	for _, d := range disks {
		targetSource, ok := targetSources[d.TargetDev]
		if !ok || targetSource == "" {
			continue
		}
		diskXML := buildDiskDeviceXML(d, targetSource)
		if err := domain.Domain.UpdateDeviceFlags(diskXML, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG); err != nil {
			return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to update disk %s: %w", d.TargetDev, err))
		}
		updated = true
	}

	if !updated {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("no backingStore entries found to restore"))
	}

	return nil
}
