package external

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func ListExternalSnapshots(domain *domCon.Domain) ([]string, error) {
	if domain == nil || domain.Domain == nil {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	snaps, err := domain.Domain.ListAllSnapshots(0)
	if err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list snapshots: %w", err))
	}

	names := make([]string, 0, len(snaps))
	for _, s := range snaps {
		isExternal, err := isExternalSnapshot(&s)
		if err == nil && isExternal {
			name, err := s.GetName()
			if err == nil {
				names = append(names, name)
			}
		}
		s.Free()
	}

	return names, nil
}
