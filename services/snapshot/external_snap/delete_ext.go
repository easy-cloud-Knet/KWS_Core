package external

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"libvirt.org/go/libvirt"
)

func DeleteExternalSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}
	if snapName == "" {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot name required"))
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

	if err := target.Delete(0); err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to delete external snapshot %s: %w", snapName, err))
	}

	return nil
}
