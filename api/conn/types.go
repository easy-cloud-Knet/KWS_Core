package conn

import (
	"sync"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
	"libvirt.org/go/libvirt"
)

type LibvirtInst = *libvirt.Connect

type Domain struct{
	Domain *libvirt.Domain
	DomainMutex sync.Mutex
}

type DomainDeviceManager struct{	
}
// managing attachable devices for vm, vcpu,internet interface ... 
type DomainStatusManager struct{
	DomainState libvirt.DomainState
	UUID string
}
// managing domain status, deleting, shutting down.... 
// need to add advanced feature like updating state, 
// or setting call back for state update

type  DomainControl interface{
	DomainStatus()
}


///////////////////////////////////////////////////

type DomainDataType uint

const (
	PowerStaus DomainDataType =iota //0
	BasicInfo	//1 .... 
	GuestInfoUser // 2
	GuestInfoOS// 3
	GuestInfoFS//4  이 숫자랑만 맞춰주면 됨 
	GuestInfoDisk
)	
type DomainDeleteType uint 

const (
	HardDelete DomainDeleteType =iota
	SoftDelete 	
)


type DomainState struct{
	DomainState libvirt.DomainState `json:"currentState"`
	//type reference 참고
	UUID string `json:"UUID"`
	Users []libvirt.DomainGuestInfoUser `json:"Guest Info"`
}

type DomainInfo struct{
	State libvirt.DomainState `json:"state"`
	MaxMem uint64 `json:"maxmem"`
	Memory uint64 `json:"memory"`
	NrVirtCpu uint `json:"nrVirtCpu"`
	CpuTime uint64 `json:"cpuTime"`
}
type DataTypeHandler interface{
	GetInfo(*Domain) error
}

////////////////////////interface uniformed function for various infoType 

type DomainDetail struct{
	DataHandle []DataTypeHandler 
	DomainSeeker DomainSeeker
}

type DomainController struct{
	DomainSeeker *DomainSeekingByUUID
}

type DomainGeneratorLocal struct{
	DomainStatusManager *DomainStatusManager
	DataParsor parsor.DomainGenerator
	OS string
}

type DomainTerminator struct{
	DomainSeeker DomainSeeker
}
type DomainDeleter struct{
	DomainSeeker DomainSeeker
	DomainStatusManager *DomainStatusManager
	DeletionType DomainDeleteType
}

////////////////////////////////////////////////////service structures for 
/////// functnions like delete, generate, stateControl 

type DomainSeekingByUUID struct{
	LibvirtInst *libvirt.Connect
	UUID string 
	Domain []*Domain
}

type DomainSeekingByStatus struct{
	LibvirtInst *libvirt.Connect
	Status libvirt.ConnectListAllDomainsFlags 
	DomList []*Domain
}

type DomainSeeker interface{
	SetDomain() (error)
	ReturnDomain()([]*Domain,error)
}
/////////////////////////////interface seeking certain Domain







 