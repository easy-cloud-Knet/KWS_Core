package conn

import (
	"sync"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
	"libvirt.org/go/libvirt"
)

type LibvirtInst = *libvirt.Connect

type Domain struct {
	Domain      *libvirt.Domain
	DomainMutex sync.Mutex
}
// managing attachable devices for vm, vcpu,internet interface ...
type DomainStatusManager struct {
	DomainState libvirt.DomainState
	UUID        string
}

// managing domain status, deleting, shutting down....
// need to add advanced feature like updating state,
// or setting call back for state update



///////////////////////////////////////////////////

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

type DomainDeleteType uint

const (
	HardDelete DomainDeleteType = iota
	SoftDelete
)

type DomainState struct {
	DomainState libvirt.DomainState `json:"currentState"`
	UUID  string                        `json:"UUID"`
	Users []libvirt.DomainGuestInfoUser `json:"Guest Info"`
}

type DomainInfo struct {
	State     libvirt.DomainState `json:"state"`
	MaxMem    uint64              `json:"maxmem"`
	Memory    uint64              `json:"memory"`
	NrVirtCpu uint                `json:"nrVirtCpu"`
	CpuTime   uint64              `json:"cpuTime"`
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

type SystemInfo struct {
	Memory MemoryInfo `json:"memory"`
	Disks  DiskInfo   `json:"disks"`
}

//

type DataTypeHandler interface {
	GetInfo(*Domain) error
}

type DomainDetail struct {
	DataHandle   []DataTypeHandler
	DomainSeeker DomainSeeker
}
////////////////////////interface uniformed function for various infoType

type DomainGeneratorLocal struct {
	DomainStatusManager *DomainStatusManager
	DataParsor          parsor.DomainGenerator
	OS                  string
}

type DomainTerminator struct {
	DomainSeeker DomainSeeker
}
type DomainDeleter struct {
	DomainSeeker        DomainSeeker
	DomainStatusManager *DomainStatusManager
	DeletionType        DomainDeleteType
}

////////////////////////////////////////////////////service structures for
/////// functnions like delete, generate, stateControl

type DomainSeekingByUUID struct {
	LibvirtInst *libvirt.Connect
	UUID        string
	Domain      []*Domain
}

type DomainSeekingByStatus struct {
	LibvirtInst *libvirt.Connect
	Status      libvirt.ConnectListAllDomainsFlags
	DomList     []*Domain
}

type DomainSeeker interface {
	SetDomain() error
	ReturnDomain() ([]*Domain, error)
}

/////////////////////////////interface seeking certain Domain

//hostinfo_

type HostDetail struct {
	HostDataHandle []HostDataTypeHandler
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
