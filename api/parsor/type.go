package parsor

import (
	"encoding/xml"
)

type VM_CREATE_XML struct {
	XMLName   xml.Name   `xml:"domain"`
	Type      string     `xml:"type,attr"`
	Name      string     `xml:"name"`
	UUID      string     `xml:"uuid"`
	Memory    Memory     `xml:"memory"`
	VCPU      VCPU       `xml:"vcpu"`
	Features  Features   `xml:"features"`
	OS        OS         `xml:"os"`
	Devices   Devices    `xml:"devices"`
}

type Memory struct {
	Unit string `xml:"unit,attr"`
	Size int    `xml:",chardata"`
}

type VCPU struct {
	Placement string `xml:"placement,attr"`
	Count     int    `xml:",chardata"`
}

type Features struct {
	ACPI ACPI `xml:"acpi"`
}

type ACPI struct{}

type OS struct {
	Type OSType `xml:"type"`
	Boot Boot   `xml:"boot"`
}

type OSType struct {
	Arch string `xml:"arch,attr"`
	Type string `xml:",chardata"`
}

type Boot struct {
	Dev string `xml:"dev,attr"`
}

type Devices struct {
	Emulator   string      `xml:"emulator"`
	Disks      []Disk      `xml:"disk"`
	Serial     Serial      `xml:"serial"`
	Console    Console     `xml:"console"`
	Interfaces []Interface `xml:"interface"`
	Graphics   Graphics    `xml:"graphics"`
}

type Disk struct {
	Type   string  `xml:"type,attr"`
	Device string  `xml:"device,attr"`
	Driver Driver  `xml:"driver"`
	Source Source  `xml:"source"`
	Target Target  `xml:"target"`
	ReadOnly *ReadOnly `xml:"readonly,omitempty"`
}

type Driver struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Source struct {
	File string `xml:"file,attr"`
}

type Target struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}

type ReadOnly struct{}

type Serial struct {
	Type   string `xml:"type,attr"`
	Target TargetPort `xml:"target"`
}

type TargetPort struct {
	Port int `xml:"port,attr"`
}

type Console struct {
	Type   string `xml:"type,attr"`
	Target ConsoleTarget `xml:"target"`
}

type ConsoleTarget struct {
	Type string `xml:"type,attr"`
	Port int    `xml:"port,attr"`
}

type Interface struct {
	Type   string      `xml:"type,attr"`
	Source NetworkSource `xml:"source"`
	Model  InterfaceModel `xml:"model"`
}

type NetworkSource struct {
	Network string `xml:"network,attr"`
}

type InterfaceModel struct {
	Type string `xml:"type,attr"`
}

type Graphics struct {
	Type     string `xml:"type,attr"`
	Port     int    `xml:"port,attr"`
	AutoPort string `xml:"autoport,attr"`
}
