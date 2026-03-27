package external

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"libvirt.org/go/libvirt"
)

func DeleteExternalSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return fmt.Errorf("nil domain")
	}
	if snapName == "" {
		return fmt.Errorf("snapshot name required")
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

	if err := target.Delete(0); err != nil {
		return fmt.Errorf("failed to delete external snapshot %s: %w", snapName, err)
	}
	if err := deleteSnapshotMetadataByName(domain, snapName); err != nil {
		return fmt.Errorf("snapshot deleted but failed to update metadata: %w", err)
	}

	return nil
}
