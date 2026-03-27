package external

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
)

func ListExternalSnapshots(domain *domCon.Domain) ([]string, error) {
	if domain == nil || domain.Domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	snaps, err := domain.Domain.ListAllSnapshots(0)
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
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
