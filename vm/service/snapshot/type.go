package snapshot

import "encoding/xml"

// ============================================================================
// Snapshot Options
// ============================================================================

// SnapshotOptions defines options for creating internal snapshots.
// Internal snapshots store the snapshot data within the original disk image.
type SnapshotOptions struct {
	Description string // snapshot name
	Quiesce     bool   // Whether to quiesce the guest filesystem before snapshot
}

// ExternalSnapshotOptions defines options for creating external snapshots.
// External snapshots store the snapshot data in separate files.
type ExternalSnapshotOptions struct {
	BaseDir     string // Base directory for snapshot storage
	Description string // Human-readable description of the snapshot
	Quiesce     bool   // Whether to quiesce the guest filesystem before snapshot
	Live        bool   // Whether to create a live snapshot (for running domains)
}

// ============================================================================
// Domain XML Parsing Types (for reading disk information)
// ============================================================================
// These types are used to unmarshal domain XML and extract disk information
// needed for external snapshot operations.

// domainXML represents the root domain XML structure.
type domainXML struct {
	Devices domainDevices `xml:"devices"`
}

// domainDevices contains all device definitions in the domain.
type domainDevices struct {
	Disks []domainDisk `xml:"disk"`
}

// domainDisk represents a single disk device in the domain.
type domainDisk struct {
	Device       string              `xml:"device,attr"`
	Type         string              `xml:"type,attr"`
	Driver       *domainDriver       `xml:"driver"`
	Source       *domainSource       `xml:"source"`
	Target       *domainTarget       `xml:"target"`
	BackingStore *domainBackingStore `xml:"backingStore"`
}

// domainDriver represents disk driver configuration.
type domainDriver struct {
	Name string `xml:"name,attr"` // Driver name (e.g., "qemu")
	Type string `xml:"type,attr"` // Format type (e.g., "qcow2", "raw")
}

// domainSource represents the source file path for the disk.
type domainSource struct {
	File string `xml:"file,attr"`
}

// domainTarget represents the target device configuration.
type domainTarget struct {
	Dev string `xml:"dev,attr"` // Device name (e.g., "vda", "sda")
	Bus string `xml:"bus,attr"` // Bus type (e.g., "virtio", "scsi")
}

// domainBackingStore represents the backing file information.
type domainBackingStore struct {
	Source *domainSource `xml:"source"`
}

// diskInfo is an internal structure holding parsed disk information.
type diskInfo struct {
	TargetDev     string // Target device name
	TargetBus     string // Target bus type
	Source        string // Current source file path
	BackingSource string // Backing file path (if exists)
	Driver        string // Driver format type
	DriverName    string // Driver name
}

// ============================================================================
// Snapshot XML Generation Types (for creating snapshots)
// ============================================================================
// These types are used to marshal snapshot XML when creating snapshots.

// snapshotXML represents the root snapshot XML structure.
type snapshotXML struct {
	XMLName     xml.Name      `xml:"domainsnapshot"`
	Name        string        `xml:"name"`
	Description string        `xml:"description,omitempty"`
	Disks       snapshotDisks `xml:"disks"`
}

// snapshotDisks contains all disk specifications for the snapshot.
type snapshotDisks struct {
	Disks []snapshotDisk `xml:"disk"`
}

// snapshotDisk represents a single disk configuration in the snapshot.
type snapshotDisk struct {
	Name     string          `xml:"name,attr"`        // Disk target device name
	Snapshot string          `xml:"snapshot,attr"`    // Snapshot mode ("external", "internal", "no")
	Driver   *snapshotDriver `xml:"driver,omitempty"` // Optional driver configuration
	Source   *snapshotSource `xml:"source,omitempty"` // Optional source file path
}

// snapshotDriver represents driver configuration for snapshot disk.
type snapshotDriver struct {
	Type string `xml:"type,attr"` // Format type (e.g., "qcow2")
}

// snapshotSource represents the snapshot file path.
type snapshotSource struct {
	File string `xml:"file,attr"`
}
