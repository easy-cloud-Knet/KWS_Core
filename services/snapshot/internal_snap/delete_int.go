package internal

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func DeleteSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	return deleteSnapshot(newInternalSnapshotDomain(domain.Domain), snapName)
}

func deleteSnapshot(domain snapshotDomain, snapName string) error {
	if domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	snaps, err := domain.ListAllSnapshots()
	if err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list snapshots: %w", err))
	}
	defer freeSnapshotHandles(snaps)

	target := findSnapshotByName(snaps, snapName)
	if target != nil {
		if err := target.Delete(); err != nil {
			return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to delete snapshot %s: %w", snapName, err))
		}
		return nil
	}

	return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("snapshot %s not found", snapName))
}
