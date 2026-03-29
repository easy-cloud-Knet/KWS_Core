package internal

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func ListSnapshots(domain *domCon.Domain) ([]string, error) {
	if domain == nil || domain.Domain == nil {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	return listSnapshots(newInternalSnapshotDomain(domain.Domain))
}

func listSnapshots(domain internalSnapshotDomain) ([]string, error) {
	if domain == nil {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	snaps, err := domain.ListAllSnapshots()
	if err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list snapshots: %w", err))
	}

	names := make([]string, 0, len(snaps))
	for _, s := range snaps {
		n, err := s.Name()
		if err == nil {
			names = append(names, n)
		}
		s.Free()
	}

	return names, nil
}
