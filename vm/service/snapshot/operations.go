package snapshot

import (
	"fmt"
	"time"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"libvirt.org/go/libvirt"
)

// CreateSnapshot creates a libvirt snapshot and records basic metadata.
func CreateSnapshot(domain *domCon.Domain, name string) (string, error) {
	if domain == nil || domain.Domain == nil {
		return "", fmt.Errorf("nil domain")
	}

	snapXML := fmt.Sprintf(`<domainsnapshot><name>%s</name><description>snapshot created by KWS</description></domainsnapshot>`, name)

	snap, err := domain.Domain.CreateSnapshotXML(snapXML, 0)
	if err != nil {
		return "", fmt.Errorf("failed to create snapshot: %w", err)
	}
	defer snap.Free()

	snapName, err := snap.GetName()
	if err != nil {
		return "", fmt.Errorf("snapshot created but failed to read name: %w", err)
	}

	// try to collect additional snapshot info from libvirt
	xmlDesc, _ := snap.GetXMLDesc(0)
	// persist metadata (include libvirt name and xml if available)
	meta := SnapshotMeta{
		Name:        snapName,
		LibvirtName: snapName,
		XMLDesc:     xmlDesc,
		CreatedAt:   time.Now(),
	}
	_ = SaveSnapshotMetaForDomain(domain, meta)

	return snapName, nil
}

// ListSnapshots lists snapshot names for the domain.
func ListSnapshots(domain *domCon.Domain) ([]string, error) {
	if domain == nil || domain.Domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	snaps, err := domain.Domain.ListAllSnapshots(0)
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
	}

	names := make([]string, 0, len(snaps))
	for _, s := range snaps {
		n, err := s.GetName()
		if err == nil {
			names = append(names, n)
		}
		s.Free()
	}

	return names, nil
}

// RevertToSnapshot reverts the domain to the given snapshot name.
func RevertToSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return fmt.Errorf("nil domain")
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
		return fmt.Errorf("snapshot %s not found", snapName)
	}

	// Call the libvirt-go binding's DomainSnapshot.RevertToSnapshot method.
	// Use zero flags for default behavior.
	if err := target.RevertToSnapshot(0); err != nil {
		return fmt.Errorf("failed to revert to snapshot %s: %w", snapName, err)
	}

	return nil
}
