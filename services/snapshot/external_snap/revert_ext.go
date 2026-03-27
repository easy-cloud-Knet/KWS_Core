package external

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"libvirt.org/go/libvirt"
)

func RevertExternalSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return fmt.Errorf("nil domain")
	}

	active, err := domain.Domain.IsActive()
	if err == nil && active {
		return fmt.Errorf("external snapshot revert requires the domain to be shut down")
	}

	snaps, err := domain.Domain.ListAllSnapshots(0)
	if err != nil {
		return fmt.Errorf("failed to list snapshots: %w", err)
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
			return err
		}
		if !isExternal {
			return fmt.Errorf("snapshot %s is not external", snapName)
		}
		target = &snaps[i]
		break
	}

	if target == nil {
		return fmt.Errorf("snapshot %s not found", snapName)
	}

	targetSources, err := extractExternalSnapshotSources(target)
	if err != nil {
		return err
	}
	if len(targetSources) == 0 {
		return fmt.Errorf("snapshot %s has no external disk sources", snapName)
	}

	disks, err := listFileDisks(domain)
	if err != nil {
		return err
	}

	updated := false
	for _, d := range disks {
		targetSource, ok := targetSources[d.TargetDev]
		if !ok || targetSource == "" {
			continue
		}
		diskXML := buildDiskDeviceXML(d, targetSource)
		if err := domain.Domain.UpdateDeviceFlags(diskXML, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG); err != nil {
			return fmt.Errorf("failed to update disk %s: %w", d.TargetDev, err)
		}
		updated = true
	}

	if !updated {
		return fmt.Errorf("no backingStore entries found to restore")
	}

	return nil
}
