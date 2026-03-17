package snapshot

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"libvirt.org/go/libvirt"
)

// ============================================================================
// External Snapshot Management
// ============================================================================
// External snapshots store snapshot data in separate files, allowing for
// more flexible backup and recovery strategies. The original disk becomes
// read-only and a new overlay file is created to track changes.

// CreateExternalSnapshot creates an external snapshot for the specified domain.
// External snapshots store the snapshot data in separate files instead of
// within the original disk image, making them suitable for live backups and
// more complex snapshot chains.
//
// Parameters:
//   - domain: The domain for which to create the snapshot
//   - name: The name of the snapshot (must be safe for filesystem use)
//   - opts: Optional settings including description, quiesce, and live flags
//
// Returns the created snapshot name or an error.
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

// ListExternalSnapshots returns all external snapshot names for the domain.
// Only snapshots that have at least one disk marked as "external" are included.
//
// Parameters:
//   - domain: The domain to list snapshots for
//
// Returns a slice of external snapshot names or an error.
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

// DeleteExternalSnapshot removes metadata for an external snapshot name from
// libvirt and the local external snapshot metadata store.
func DeleteExternalSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return fmt.Errorf("nil domain")
	}
	if snapName == "" {
		return fmt.Errorf("snapshot name required")
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
		name, err := snaps[i].GetName()
		if err != nil || name != snapName {
			continue
		}
		isExternal, err := isExternalSnapshot(&snaps[i])
		if err != nil {
			return err
		}
		if !isExternal {
			return fmt.Errorf("snapshot %s is not external", snapName)
		}
		target = &snaps[i]
		break
	}

	if target == nil {
		return fmt.Errorf("snapshot %s not found", snapName)
	}

	if err := target.Delete(0); err != nil {
		return fmt.Errorf("failed to delete external snapshot %s: %w", snapName, err)
	}

	if err := deleteSnapshotMetadataByName(domain, snapName); err != nil {
		return fmt.Errorf("snapshot deleted but failed to update metadata: %w", err)
	}

	return nil
}

// 외부 스냅샷 삭제 기능 구현해야함.
//!!!!!!!!!!

// RevertExternalSnapshot reverts the domain to a previously created external snapshot.
// The domain must be shut down before reverting to an external snapshot.
// This operation updates the domain's disk configuration to point back to the
// snapshot's backing files.
//
// Parameters:
//   - domain: The domain to revert
//   - snapName: The name of the external snapshot to revert to
//
// Returns an error if the operation fails or if the domain is running.
func RevertExternalSnapshot(domain *domCon.Domain, snapName string) error {
	if domain == nil || domain.Domain == nil {
		return fmt.Errorf("nil domain")
	}

	active, err := domain.Domain.IsActive()
	if err == nil && active {
		return fmt.Errorf("external snapshot revert requires the domain to be shut down")
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
		name, err := snaps[i].GetName()
		if err != nil || name != snapName {
			continue
		}
		isExternal, err := isExternalSnapshot(&snaps[i])
		if err != nil {
			return err
		}
		if !isExternal {
			return fmt.Errorf("snapshot %s is not external", snapName)
		}
		target = &snaps[i]
		break
	}

	if target == nil {
		return fmt.Errorf("snapshot %s not found", snapName)
	}

	disks, err := listFileDisks(domain)
	if err != nil {
		return err
	}

	updated := false
	for _, d := range disks {
		if d.BackingSource == "" {
			continue
		}
		diskXML := buildDiskDeviceXML(d, d.BackingSource)
		if err := domain.Domain.UpdateDeviceFlags(diskXML, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG); err != nil {
			return fmt.Errorf("failed to update disk %s: %w", d.TargetDev, err)
		}
		updated = true
	}

	if !updated {
		return fmt.Errorf("no backingStore entries found to restore")
	}

	return nil
}

// MergeExternalSnapshot merges the current external overlay into its immediate
// backing file for one or more disks.
//
// This operation is supported only when the domain is running.
// If targetDisk is empty, all eligible file-based disks are merged.
//
// Parameters:
//   - domain: The domain whose external snapshot layers will be merged
//   - targetDisk: Optional target disk device name (e.g., "vda")
//
// Returns the list of merged disk targets or an error.
func MergeExternalSnapshot(domain *domCon.Domain, targetDisk string) ([]string, error) {
	if domain == nil || domain.Domain == nil {
		return nil, fmt.Errorf("nil domain")
	}

	active, err := domain.Domain.IsActive()
	if err != nil {
		return nil, fmt.Errorf("failed to check domain state: %w", err)
	}
	if !active {
		return nil, fmt.Errorf("external snapshot merge requires the domain to be running")
	}

	disks, err := listFileDisks(domain)
	if err != nil {
		return nil, err
	}

	store, _, err := loadExternalSnapshotMetadata(domain)
	if err != nil {
		return nil, err
	}

	merged := make([]string, 0, len(disks))
	mergedMetaKeys := make(map[string]struct{})
	for _, d := range disks {
		if targetDisk != "" && d.TargetDev != targetDisk {
			continue
		}

		backingSource := d.BackingSource
		if backingSource == "" {
			resolvedBacking, ok := findDiskBackingFromMetadata(store, d.TargetDev, d.Source)
			if ok {
				backingSource = resolvedBacking
			}
		}

		if backingSource == "" {
			continue
		}

		flags := libvirt.DOMAIN_BLOCK_COMMIT_ACTIVE | libvirt.DOMAIN_BLOCK_COMMIT_DELETE
		if err := domain.Domain.BlockCommit(d.TargetDev, backingSource, d.Source, 0, flags); err != nil {
			return nil, fmt.Errorf("failed to start block commit on disk %s: %w", d.TargetDev, err)
		}

		if err := waitBlockJobReady(domain.Domain, d.TargetDev, 2*time.Minute); err != nil {
			return nil, err
		}

		if err := domain.Domain.BlockJobAbort(d.TargetDev, libvirt.DOMAIN_BLOCK_JOB_ABORT_PIVOT); err != nil {
			return nil, fmt.Errorf("failed to pivot disk %s after commit: %w", d.TargetDev, err)
		}

		diskXML := buildDiskDeviceXML(d, backingSource)
		if err := domain.Domain.UpdateDeviceFlags(diskXML, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG); err != nil {
			return nil, fmt.Errorf("failed to update disk %s after merge: %w", d.TargetDev, err)
		}

		mergedMetaKeys[d.TargetDev+"|"+d.Source] = struct{}{}
		merged = append(merged, d.TargetDev)
	}

	if targetDisk != "" && len(merged) == 0 {
		return nil, fmt.Errorf("disk %s is not mergeable or has no external backing chain", targetDisk)
	}
	if len(merged) == 0 {
		return nil, fmt.Errorf("no mergeable external snapshot disks found")
	}

	if err := pruneMergedDisksFromMetadata(domain, mergedMetaKeys); err != nil {
		return nil, fmt.Errorf("merge completed but failed to update metadata: %w", err)
	}

	return merged, nil
}

// waitBlockJobReady waits until a block commit job has copied all data and is
// ready for pivot.
func waitBlockJobReady(domain *libvirt.Domain, disk string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		job, err := domain.GetBlockJobInfo(disk, 0)
		if err != nil {
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for block job on disk %s: %w", disk, err)
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if job.End > 0 && job.Cur >= job.End {
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for block commit to complete on disk %s", disk)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// listFileDisks retrieves all file-based disk devices from the domain.
// Only disks with device type "disk" and type "file" are included.
//
// Returns a slice of diskInfo structures containing disk configuration details.
func listFileDisks(domain *domCon.Domain) ([]diskInfo, error) {
	xmlDesc, err := domain.Domain.GetXMLDesc(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain xml: %w", err)
	}

	var doc domainXML
	if err := xml.Unmarshal([]byte(xmlDesc), &doc); err != nil {
		return nil, fmt.Errorf("failed to parse domain xml: %w", err)
	}

	out := make([]diskInfo, 0, len(doc.Devices.Disks))
	for _, d := range doc.Devices.Disks {
		if d.Device != "disk" || d.Type != "file" {
			continue
		}
		if d.Source == nil || d.Target == nil || d.Target.Dev == "" || d.Source.File == "" {
			continue
		}
		driverType := ""
		driverName := ""
		if d.Driver != nil {
			driverType = d.Driver.Type
			driverName = d.Driver.Name
		}
		backingSource := ""
		if d.BackingStore != nil && d.BackingStore.Source != nil {
			backingSource = d.BackingStore.Source.File
		}
		out = append(out, diskInfo{
			TargetDev:     d.Target.Dev,
			TargetBus:     d.Target.Bus,
			Source:        d.Source.File,
			BackingSource: backingSource,
			Driver:        driverType,
			DriverName:    driverName,
		})
	}

	return out, nil
}

// buildDiskDeviceXML constructs XML for a disk device configuration.
// Used when reverting external snapshots to update disk sources.
//
// Parameters:
//   - info: The disk information structure
//   - source: The source file path to use for the disk
//
// Returns the XML string for the disk device.
func buildDiskDeviceXML(info diskInfo, source string) string {
	driverXML := ""
	if info.Driver != "" || info.DriverName != "" {
		driverXML = "<driver"
		if info.DriverName != "" {
			driverXML += fmt.Sprintf(" name='%s'", info.DriverName)
		}
		if info.Driver != "" {
			driverXML += fmt.Sprintf(" type='%s'", info.Driver)
		}
		driverXML += "/>"
	}

	targetXML := fmt.Sprintf("<target dev='%s'", info.TargetDev)
	if info.TargetBus != "" {
		targetXML += fmt.Sprintf(" bus='%s'", info.TargetBus)
	}
	targetXML += "/>"

	return fmt.Sprintf("<disk type='file' device='disk'>%s<source file='%s'/>%s</disk>", driverXML, source, targetXML)
}

// isExternalSnapshot determines if a snapshot is an external snapshot.
// A snapshot is considered external if at least one disk has snapshot="external".
//
// Returns true if the snapshot is external, false otherwise.
func isExternalSnapshot(snapshot *libvirt.DomainSnapshot) (bool, error) {
	if snapshot == nil {
		return false, fmt.Errorf("nil snapshot")
	}

	xmlDesc, err := snapshot.GetXMLDesc(0)
	if err != nil {
		return false, fmt.Errorf("failed to get snapshot xml: %w", err)
	}

	var doc snapshotXML
	if err := xml.Unmarshal([]byte(xmlDesc), &doc); err != nil {
		return false, fmt.Errorf("failed to parse snapshot xml: %w", err)
	}

	for _, d := range doc.Disks.Disks {
		if strings.EqualFold(d.Snapshot, "external") {
			return true, nil
		}
	}

	return false, nil
}

// isSafeSnapshotName validates that a snapshot name is safe for filesystem use.
// The name must not contain path separators, parent directory references (..)
// or result in a different name after path cleaning.
//
// Returns true if the name is safe, false otherwise.
func isSafeSnapshotName(name string) bool {
	if name == "" {
		return false
	}
	clean := filepath.Clean(name)
	if clean != name {
		return false
	}
	if strings.Contains(name, "..") {
		return false
	}
	if strings.ContainsAny(name, `/\\`) {
		return false
	}
	return true
}
