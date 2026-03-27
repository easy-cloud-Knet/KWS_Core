package internal

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"libvirt.org/go/libvirt"
)

func DeleteSnapshot(domain *domCon.Domain, snapName string) error {
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

	for _, s := range snaps {
		n, err := s.GetName()
		if err != nil {
			continue
		}
		if n == snapName {
			if err := s.Delete(0); err != nil {
				return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to delete snapshot %s: %w", snapName, err))
			}
			return nil
		}
	}

	return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("snapshot %s not found", snapName))
}
