package snapshot

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"libvirt.org/go/libvirt"
)

type domainXML struct {
	Devices domainDevices `xml:"devices"`
}

type domainDevices struct {
	Disks []domainDisk `xml:"disk"`
}

type domainDisk struct {
	Device       string              `xml:"device,attr"`
	Type         string              `xml:"type,attr"`
	Driver       *domainDriver       `xml:"driver"`
	Source       *domainSource       `xml:"source"`
	Target       *domainTarget       `xml:"target"`
	BackingStore *domainBackingStore `xml:"backingStore"`
}

type domainDriver struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type domainSource struct {
	File string `xml:"file,attr"`
}

type domainTarget struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}

type domainBackingStore struct {
	Source *domainSource `xml:"source"`
}

type diskInfo struct {
	TargetDev     string
	TargetBus     string
	Source        string
	BackingSource string
	Driver        string
	DriverName    string
}

type snapshotXML struct {
	XMLName     xml.Name      `xml:"domainsnapshot"`
	Name        string        `xml:"name"`
	Description string        `xml:"description,omitempty"`
	Disks       snapshotDisks `xml:"disks"`
}

type snapshotDisks struct {
	Disks []snapshotDisk `xml:"disk"`
}

type snapshotDisk struct {
	Name     string          `xml:"name,attr"`
	Snapshot string          `xml:"snapshot,attr"`
	Driver   *snapshotDriver `xml:"driver,omitempty"`
	Source   *snapshotSource `xml:"source,omitempty"`
}

type snapshotDriver struct {
	Type string `xml:"type,attr"`
}

type snapshotSource struct {
	File string `xml:"file,attr"`
}

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

	snapDisks := make([]snapshotDisk, 0, len(disks))
	for _, d := range disks {
		var driver *snapshotDriver
		if d.Driver != "" {
			driver = &snapshotDriver{Type: d.Driver}
		}

		snapDisks = append(snapDisks, snapshotDisk{
			Name:     d.TargetDev,
			Snapshot: "external",
			Driver:   driver,
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

	return snapName, nil
}

// ListExternalSnapshots lists only external snapshot names for the domain.
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

// RevertExternalSnapshot reverts the domain to an external snapshot.
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
