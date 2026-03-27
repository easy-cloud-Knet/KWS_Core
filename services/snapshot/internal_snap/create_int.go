package internal

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"libvirt.org/go/libvirt"
)

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
