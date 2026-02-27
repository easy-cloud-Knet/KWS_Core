package snapshot

import "encoding/xml"

// External snapshot XML types (used by external.go).
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

// Internal snapshot options (used by operations.go).
type SnapshotOptions struct {
	Description string
	Quiesce     bool
}
