package snapshot

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"libvirt.org/go/libvirt"
)

// ============================================================================
// Internal Snapshot Management
// ============================================================================
// Internal snapshots store snapshot data within the original disk image file.
// They are simpler to manage but require qcow2 or similar formats that support
// internal snapshots. All snapshot data is contained in a single file.

// CreateSnapshot creates an internal snapshot for the specified domain.
// Internal snapshots are stored within the original disk image file and are
// suitable for qcow2 and other formats that support internal snapshot storage.
//
// Parameters:
//   - domain: The domain for which to create the snapshot
//   - name: The name of the snapshot
//   - opts: Optional settings including description and quiesce flag
//
// Returns the created snapshot name or an error.
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

// ListSnapshots returns all snapshot names for the specified domain.
// This includes both internal and external snapshots.
//
// Parameters:
//   - domain: The domain to list snapshots for
//
// Returns a slice of snapshot names or an error.
func ListSnapshots(domain *domCon.Domain) ([]string, error) {
	if domain == nil || domain.Domain == nil {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	var listFlags libvirt.DomainSnapshotListFlags
	snaps, err := domain.Domain.ListAllSnapshots(listFlags)
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

// RevertToSnapshot reverts the domain to a previously created snapshot.
// The domain will be restored to the exact state it was in when the snapshot
// was created, including memory state if the snapshot was taken while running.
//
// Parameters:
//   - domain: The domain to revert
//   - snapName: The name of the snapshot to revert to
//
// Returns an error if the snapshot doesn't exist or revert fails.
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

// DeleteSnapshot permanently deletes a snapshot by name.
// This removes the snapshot metadata and, depending on the snapshot type,
// may merge snapshot data back into the base disk.
//
// Parameters:
//   - domain: The domain containing the snapshot
//   - snapName: The name of the snapshot to delete
//
// Returns an error if the snapshot doesn't exist or deletion fails.
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
