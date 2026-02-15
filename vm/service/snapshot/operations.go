package snapshot

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"libvirt.org/go/libvirt"
)

type SnapshotOptions struct {
	Description string
	Quiesce     bool
}

// CreateSnapshot creates a libvirt snapshot and records basic metadata.
func CreateSnapshot(domain *domCon.Domain, name string, opts *SnapshotOptions) (string, error) {
	if domain == nil || domain.Domain == nil {
		return "", virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	description := "snapshot created by KWS"
	if opts != nil && opts.Description != "" {
		description = opts.Description
	}

	snapXML := fmt.Sprintf(`<domainsnapshot><name>%s</name><description>%s</description></domainsnapshot>`, name, description)

	flags := libvirt.DomainSnapshotCreateFlags(0)
	if opts != nil && opts.Quiesce {
		flags |= libvirt.DOMAIN_SNAPSHOT_CREATE_QUIESCE
	}

	snap, err := domain.Domain.CreateSnapshotXML(snapXML, flags)
	if err != nil {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to create snapshot: %w", err))
	}
	defer snap.Free()

	snapName, err := snap.GetName()
	if err != nil {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("snapshot created but failed to read name: %w", err))
	}

	return snapName, nil
}

// ListSnapshots lists snapshot names for the domain.
func ListSnapshots(domain *domCon.Domain) ([]string, error) {
	if domain == nil || domain.Domain == nil {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	snaps, err := domain.Domain.ListAllSnapshots(0)
	if err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list snapshots: %w", err))
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
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
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

// DeleteSnapshot deletes a snapshot by name.
func DeleteSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
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
