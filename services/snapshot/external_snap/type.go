package external

import "encoding/xml"

type ExternalSnapshotOptions struct {
	BaseDir     string
	Description string
}

// QemuImg abstracts qemu-img command execution for testability.
type QemuImg interface {
	Create(backingFile, backingFormat, overlayPath string) error
	Info(diskPath string) (backingFile, backingFormat string, err error)
	// Convert flattens the full backing chain of src into a new standalone dst file.
	// Unlike Commit, it never writes to any backing file so shared base images stay untouched.
	Convert(src, dst string) error
	// Commit writes all changes recorded in overlay into base, collapsing every
	// intermediate layer between them. Only base needs a write lock; files above
	// base in the chain are read-only during the operation.
	Commit(overlay, base string) error
}

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
