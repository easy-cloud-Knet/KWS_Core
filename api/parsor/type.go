package parsor

import (
	"encoding/xml"

	network "github.com/easy-cloud-Knet/KWS_Core.git/api/net"
)

//type IP []byte

type Create_VM_Method uint


type DomainGenerator struct{
	YamlParsor YamlController 
	DeviceDefiner VM_CREATE_XML
}
// struct for detailed generation config

type YamlController interface{
	Parse_data(*VM_Init_Info)  error
	FileConfig(string) error 
}

const (
	CREATE_WITH_XML Create_VM_Method = iota+1
	type1
	type2
	type3
)

type Meta_data_yaml struct{
	Instance_ID string `yaml:"instance-id"`
	Local_Host_Id string `yaml:"local-hostname"`
}

type User_specific struct{
	Name string `yaml:"name,omitempty"`
	Passwd string `yaml:"passwd,omitempty"`
	Lock_passwd bool `yaml:"lock_passwd"`
	Ssh_authorized_keys []string `yaml:"ssh_authorized_keys,omitempty"`
	Groups string `yaml:"groups,omitempty"`
	SuGroup string `yaml:"sudo,omitempty"`
	Shell string ` yaml:"shell,omitempty"`
}

type User_write_file struct{
	Path string `yaml:"path"`
	Permissions string `yaml:"permissions"`
	Content string `yaml:"content"`
}
type User_data_yaml struct{
	PackageUpdatable bool `yaml:"package_update"`
	PredownProjects []string `yaml:"packages"`
	Users []interface{}  `yaml:"users"`
	Write_files []User_write_file `yaml:"write_files"`
	Runcmd []string `yaml:"runcmd"`
}






type User_info_VM struct {
	Name string `json:"name"`
	Groups string `json:"groups"`
	PassWord string `json:"passWord"`
}

type VM_Init_Info struct{
	DomType string `json:"domType"`
	DomName string `json:"domName"`
	UUID string `json:"uuid"`
	OS string `json:"os"`
	NetworkType string `json:"netType"`
	HardwardInfo HardwareInfo `json:"HWInfo"`
	NetConf network.NetDefine `json:"network"`
	IPs []string `json:"ips"`
	Method Create_VM_Method `json:"method"`
	Users []User_info_VM `json:"users"`
}

type HardwareInfo struct{
	CPU int `json:"cpu"`
	Memory int `json:"memory"`
}
// gonna replace fields in VM_Init_Info
//structure,need to modify parsor when implement this



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
	Channels   []Channel   `xml:"channel,omitempty"` // 새 필드 추가
}

type Channel struct {
	Type   string          `xml:"type,attr"`
	Source ChannelSource   `xml:"source"`
	Target ChannelTarget   `xml:"target"`
}

type ChannelSource struct {
	Mode string `xml:"mode,attr"`
}

type ChannelTarget struct {
	Type string `xml:"type,attr"`
	Name string `xml:"name,attr"`
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
	Bridge string `xml:"bridge,attr"`
}

type InterfaceModel struct {
	Type string `xml:"type,attr"`
}

type Graphics struct {
	Type     string `xml:"type,attr"`
	Port     int    `xml:"port,attr"`
	AutoPort string `xml:"autoport,attr"`
}
