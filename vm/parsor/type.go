package parsor

import (
	"encoding/xml"

	network "github.com/easy-cloud-Knet/KWS_Core.git/net"
)

// struct for detailed generation config

type VM_Init_Info struct {
	DomType      string            `json:"domType"`
	DomName      string            `json:"domName"`
	UUID         string            `json:"uuid"`
	OS           string            `json:"os"`
	HardwardInfo HardwareInfo      `json:"HWInfo"`
	NetConf      network.NetDefine `json:"network"`
	Users        []User_info_VM    `json:"users"`
}

type User_info_VM struct {
	Name                string   `json:"name"`
	Groups              string   `json:"groups"`
	PassWord            string   `json:"passWord"`
	Ssh_authorized_keys []string `json:"ssh"`
}

type HardwareInfo struct {
	CPU    int `json:"cpu"`
	Memory int `json:"memory"`
}



// gonna replace fields in VM_Init_Info
//structure,need to modify parsor when implement this

type VM_CREATE_XML struct {
	XMLName  xml.Name `xml:"domain"`
	Type     string   `xml:"type,attr"`
	Name     string   `xml:"name"`
	UUID     string   `xml:"uuid"`
	Memory   Memory   `xml:"memory"`
	VCPU     VCPU     `xml:"vcpu"`
	Features Features `xml:"features"`
	OS       OS       `xml:"os"`
	Devices  Devices  `xml:"devices"`
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
	Channels   []Channel   `xml:"channel,omitempty"` // 새 필드 추가
}

type Channel struct {
	Type   string        `xml:"type,attr"`
	Source ChannelSource `xml:"source"`
	Target ChannelTarget `xml:"target"`
}

type ChannelSource struct {
	Mode string `xml:"mode,attr"`
}

type ChannelTarget struct {
	Type string `xml:"type,attr"`
	Name string `xml:"name,attr"`
}

type Disk struct {
	Type     string    `xml:"type,attr"`
	Device   string    `xml:"device,attr"`
	Driver   Driver    `xml:"driver"`
	Source   Source    `xml:"source"`
	Target   Target    `xml:"target"`
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
	Type   string     `xml:"type,attr"`
	Target TargetPort `xml:"target"`
}

type TargetPort struct {
	Port int `xml:"port,attr"`
}

type Console struct {
	Type   string        `xml:"type,attr"`
	Target ConsoleTarget `xml:"target"`
}

type ConsoleTarget struct {
	Type string `xml:"type,attr"`
	Port int    `xml:"port,attr"`
}

type Interface struct {
	Type   string         `xml:"type,attr"`
	Source NetworkSource  `xml:"source"`
	Model  InterfaceModel `xml:"model"`
	Virtualport virPort  `xml:"Virtualport"`
}
type virPort struct{
	Type string `xml:"type,attr"`
}

type NetworkSource struct {
	Network string `xml:"network,attr"`
	Bridge  string `xml:"bridge,attr"`
}

type InterfaceModel struct {
	Type string `xml:"type,attr"`
}

type Graphics struct {
	Type     string `xml:"type,attr"`
	Port     int    `xml:"port,attr"`
	AutoPort string `xml:"autoport,attr"`
}
