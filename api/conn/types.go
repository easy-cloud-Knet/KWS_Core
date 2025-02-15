package conn

import (
	"sync"

	"libvirt.org/go/libvirt"
)


type DomListControl struct {
	DomainList map[string]Domain
	DomainMutex sync.Mutex 
}

type Domain struct {
	Domain      *libvirt.Domain 
}

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
	Users []libvirt.DomainGuestInfoUser `json:"Guest Info,omitempty"`
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





////////////////////////////////////////////////////////////////////////////////////
//create domain 과 관련 된 구조체들, 
// interface, Generate 로 생성 방식을 추상화 할 예졍, 
// 현재는 yaml controller 로 파일을 바로 만들어서 실행하지만, 
// user 가 입력을 하거나, 다른 파일 서버에서 받아오는 것<----- 더 빠를 수 있음
////////////////////////////////////////////////////////////////////////////////////


type DomainGenerator struct {
	Domain Domain
	DataParsor          DomainConfigGenerator
}



type DomainConfigGenerator interface {
	Generate() error 
}