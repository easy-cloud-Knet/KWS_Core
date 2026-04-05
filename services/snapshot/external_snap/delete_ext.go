package external

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func DeleteExternalSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	return deleteExternalSnapshot(newExternalSnapshotDomain(domain.Domain), snapName)
}

func deleteExternalSnapshot(domain SnapshotDomain, snapName string) error {
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

	if err := target.Delete(); err != nil {
		return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to delete external snapshot %s: %w", snapName, err))
	}

	return nil
}
