package internal

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"libvirt.org/go/libvirt"
)

func RevertToSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	var listFlags libvirt.DomainSnapshotListFlags
	snaps, err := domain.Domain.ListAllSnapshots(listFlags)
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
		n, err := snaps[i].GetName()
		if err != nil {
			continue
		}
		if n == snapName {
			target = &snaps[i]
			break
		}
	}

	if target == nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("snapshot %s not found", snapName))
	}

	if err := target.RevertToSnapshot(0); err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to revert to snapshot %s: %w", snapName, err))
	}

	return nil
}
