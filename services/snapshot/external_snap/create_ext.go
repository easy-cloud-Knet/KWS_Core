package external

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"libvirt.org/go/libvirt"
)

func CreateExternalSnapshot(domain *domCon.Domain, name string, opts *ExternalSnapshotOptions) (string, error) {
	if domain == nil || domain.Domain == nil {
		return "", fmt.Errorf("nil domain")
	}
	if name == "" {
		return "", fmt.Errorf("snapshot name required")
	}
	if !isSafeSnapshotName(name) {
		return "", fmt.Errorf("invalid snapshot name")
	}

	disks, err := listFileDisks(domain)
	if err != nil {
		return "", err
	}
	if len(disks) == 0 {
		return "", fmt.Errorf("no file-based disks found for snapshot")
	}

	snapshotRoot, err := resolveSnapshotRoot(opts)
	if err != nil {
		return "", err
	}
	domainUUID, err := resolveDomainUUID(domain)
	if err != nil {
		return "", err
	}

	snapshotDir := filepath.Join(snapshotRoot, domainUUID, "snapshots", name)
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	snapDisks := make([]snapshotDisk, 0, len(disks))
	diskMetas := make([]externalSnapshotDiskMeta, 0, len(disks))
	for _, d := range disks {
		var driver *snapshotDriver
		if d.Driver != "" {
			driver = &snapshotDriver{Type: d.Driver}
		}

		snapshotFile := filepath.Join(snapshotDir, fmt.Sprintf("%s.qcow2", d.TargetDev))
		snapDisks = append(snapDisks, snapshotDisk{
			Name:     d.TargetDev,
			Snapshot: "external",
			Driver:   driver,
			Source:   &snapshotSource{File: snapshotFile},
		})
		diskMetas = append(diskMetas, externalSnapshotDiskMeta{
			TargetDev:    d.TargetDev,
			SnapshotFile: snapshotFile,
			BackingFile:  d.Source,
		})
	}

	description := "external snapshot created by KWS"
	if opts != nil && opts.Description != "" {
		description = opts.Description
	}

	snapXML := snapshotXML{
		Name:        name,
		Description: description,
		Disks:       snapshotDisks{Disks: snapDisks},
	}

	snapBytes, err := xml.Marshal(snapXML)
	if err != nil {
		return "", fmt.Errorf("failed to build snapshot xml: %w", err)
	}

	flags := libvirt.DOMAIN_SNAPSHOT_CREATE_DISK_ONLY
	active, err := domain.Domain.IsActive()
	if opts != nil && opts.Live && err == nil && active {
		flags |= libvirt.DOMAIN_SNAPSHOT_CREATE_LIVE
	}
	if opts != nil && opts.Quiesce {
		flags |= libvirt.DOMAIN_SNAPSHOT_CREATE_QUIESCE
	}
	if len(disks) > 1 {
		flags |= libvirt.DOMAIN_SNAPSHOT_CREATE_ATOMIC
	}

	snap, err := domain.Domain.CreateSnapshotXML(string(snapBytes), flags)
	if err != nil {
		return "", fmt.Errorf("failed to create external snapshot: %w", err)
	}
	defer snap.Free()

	snapName, err := snap.GetName()
	if err != nil {
		return "", fmt.Errorf("snapshot created but failed to read name: %w", err)
	}
	if err := appendExternalSnapshotMetadata(domain, snapName, diskMetas); err != nil {
		return "", fmt.Errorf("snapshot created but failed to write metadata: %w", err)
	}

	return snapName, nil
}
