package status

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core.git/api/conn/DomCon"
	"libvirt.org/go/libvirt"
)



type HostDetail struct {
	HostDataHandle HostDataTypeHandler
}

type HostDataTypeHandler interface {
	GetHostInfo() error
}

type HostDataType uint

const (
	CpuInfo HostDataType = iota //0
	MemInfo                     //1 ....
	DiskInfoHi
	SystemInfoHi
	GeneralInfo
)

type HostCpuInfo struct {
	System float64 `json:"system_time"`
	Idle   float64 `json:"idle_time"`
	Usage  float64 `json:"usage_percent"`
}

type HostMemoryInfo struct {
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Available   uint64  `json:"available_gb"`
	UsedPercent float64 `json:"used_percent"`
}

type HostDiskInfo struct {
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Free        uint64  `json:"free_gb"`
	UsedPercent float64 `json:"used_percent"`
}

type HostSystemInfo struct {
	Uptime   uint64  `json:"uptime_seconds"`
	BootTime uint64  `json:"boot_time_epoch"`
	CPU_Temp float64 `json:"cpu_temperature,omitempty"` // no
	RAM_Temp float64 `json:"ram_temperature,omitempty"` // no
}

type HostGeneralInfo struct {
	CPU HostCpuInfo `json:"cpuInfo"`
	Memory HostMemoryInfo `json:"memoryInfo"`
	Disk HostDiskInfo `json:"DiskInfo"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Available   uint64  `json:"available_gb"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskInfo struct {
	Mountpoint  string  `json:"mountpoint"`
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Free        uint64  `json:"free_gb"`
	UsedPercent float64 `json:"used_percent"`
}


///////////////////////////////위쪽은 호스트 인포

type DomainDataType uint

const (
	DomState      DomainDataType = iota //0
	BasicInfo                           //1 ....
	GuestInfoUser                       // 2
	GuestInfoOS                         // 3
	GuestInfoFS                         //4  이 숫자랑만 맞춰주면 됨
	GuestInfoDisk
	HostInfo // -ing
)



type DomainState struct {
	DomainState libvirt.DomainState `json:"currentState"`
	UUID  string                        `json:"UUID"`
	Users []libvirt.DomainGuestInfoUser `json:"Guest Info,omitempty"`
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
	DataHandle   DataTypeHandler
	Domain *domCon.Domain
}