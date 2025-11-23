package status

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domain_status"
	"libvirt.org/go/libvirt"
)

type HostDetail struct {
	HostDataHandle HostDataTypeHandler
}

type HostDataTypeHandler interface {
	GetHostInfo(*domStatus.DomainListStatus) error
}

type HostDataType uint

const (
	CpuInfo HostDataType = iota //0
	MemInfo                     //1 ....
	DiskInfoHi
	SystemInfoHi
	GeneralInfo
	DomainallInfo
)

type HostCpuInfo struct {
	System float64 `json:"system_time"`
	Idle   float64 `json:"idle_time"`
	Usage  float64 `json:"usage_percent"`
	Desc   *domStatus.VCPUStatus `json:"vcpu_status"`
}

type HostMemoryInfo struct {
	Total          uint64  `json:"total_gb"`
	Used           uint64  `json:"used_gb"`
	Available      uint64  `json:"available_gb"`
	UsedPercent    float64 `json:"used_percent"`
	ReservedMemory uint64  `json:"reservedmem"`
	Desc   *domStatus.VCPUStatus `json:"vcpu_status"`

}

type HostDiskInfo struct {
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Free        uint64  `json:"free_gb"`
	UsedPercent float64 `json:"used_percent"`
	Desc   *domStatus.VCPUStatus `json:"vcpu_status"`

}

type HostSystemInfo struct {
	Uptime   uint64  `json:"uptime_seconds"`
	BootTime uint64  `json:"boot_time_epoch"`
	CPU_Temp float64 `json:"cpu_temperature,omitempty"` // no
	RAM_Temp float64 `json:"ram_temperature,omitempty"` // no
}

type HostGeneralInfo struct {
	CPU    HostCpuInfo    `json:"cpuInfo"`
	Memory HostMemoryInfo `json:"memoryInfo"`
	Disk   HostDiskInfo   `json:"DiskInfo"`
}

type MemoryInfo struct {
	Total          uint64  `json:"total_gb"`
	Used           uint64  `json:"used_gb"`
	Available      uint64  `json:"available_gb"`
	UsedPercent    float64 `json:"used_percent"`
	ReservedMemory uint64  `json:"reservedmem"`
}

type TotaldomainInfo struct {
	Mountpoint  string  `json:"mountpoint"`
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Free        uint64  `json:"free_gb"`
	UsedPercent float64 `json:"used_percent"`
}

///////////////////////////////위쪽은 hostinfo

type InstDataType uint

type InstDetail struct {
	AllInstDataHandle InstDataTypeHandler
}

type InstDataTypeHandler interface {
	GetAllinstInfo(LibvirtInst *libvirt.Connect) error
}

const (
	Vcpu_MaxMem InstDataType = iota //0
	//1 ....
)

type AllInstInfo struct {
	Totalmaxmem uint64 `json:"totalmaxmem"`
	TotalVCpu   uint   `json:"totalVCpu"`
}

///////////////////////////////위쪽은 allinstinfo

type DomainDataType uint

const (
	DomState      DomainDataType = iota //0
	BasicInfo                           //1 ....
	GuestInfoUser                       // 2
	GuestInfoOS                         // 3
	GuestInfoFS                         //4  이 숫자랑만 맞춰주면 됨
	GuestInfoDisk
)

type DomainState struct {
	DomainState libvirt.DomainState           `json:"currentState"`
	UUID        string                        `json:"UUID"`
	Users       []libvirt.DomainGuestInfoUser `json:"Guest Info,omitempty"`
}

type DomainInfo struct {
	State     libvirt.DomainState `json:"state"`
	MaxMem    uint64              `json:"maxmem"`
	Memory    uint64              `json:"memory"`
	NrVirtCpu uint                `json:"nrVirtCpu"`
	CpuTime   uint64              `json:"cpuTime"`
}

type DataTypeHandler interface {
	GetInfo(*domCon.Domain) error
}

type DomainDetail struct {
	DataHandle DataTypeHandler
	Domain     *domCon.Domain
}
