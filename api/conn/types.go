package conn

import (
	"sync"

	"libvirt.org/go/libvirt"
)
type DomainDataType uint



type Domain struct{
	Domain *libvirt.Domain
	DomainMutex sync.Mutex
}

type DomainDeviceManager struct{	
}
// managing attachable devices for vm, vcpu,internet interface ... 
type DomainStatusManager struct{
}
// managing domain status, deleting, shutting down, .... 

type  DomainControl interface{
	DomainStatus()
}

const (
	PowerStaus DomainDataType =iota
	BasicInfo
	GuestInfoUser
	GuestInfoOS
	GuestInfoFS
	GuestInfoDisk
)

type DataTypeHandle interface{
	Setter(DomainDataType, *libvirt.Connect)
	Getter() (DomainDataType,*libvirt.Connect)
	PowerStatus(*libvirt.Domain) 
	ReturnBasicInfo([]*libvirt.Domain) ([]*libvirt.DomainInfo,error)
	GuestInfo(*libvirt.Domain)
	
}

type InstHandler struct{
	LibvirtInst *libvirt.Connect
}

type DomainSortingByUUID[T PredefinedStructures] struct{
	TypeBase *DataType[T]
	UUID string `json:"UUID"`
}

type DomainSortingByStatus[T PredefinedStructures] struct{
	TypeBase *DataType[T]
	Status libvirt.ConnectListAllDomainsFlags `json:"Status Flag"`
}

type DataType[T PredefinedStructures] struct{
	DataType DomainDataType `json:"type"`
	LibvirtInst *libvirt.Connect
}

type DomainSeeker interface{
	returnDomain() ([]*libvirt.Domain, error)
}

// type ConnectListAllDomainsFlags uint
// const (
//	CONNECT_LIST_DOMAINS_ACTIVE         = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_ACTIVE)
// 	CONNECT_LIST_DOMAINS_INACTIVE       = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_INACTIVE)
// 	CONNECT_LIST_DOMAINS_PERSISTENT     = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_PERSISTENT)
// 	CONNECT_LIST_DOMAINS_TRANSIENT      = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_TRANSIENT)
// 	CONNECT_LIST_DOMAINS_RUNNING        = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_RUNNING)
// 	CONNECT_LIST_DOMAINS_PAUSED         = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_PAUSED)
// 	CONNECT_LIST_DOMAINS_SHUTOFF        = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_SHUTOFF)
// 	CONNECT_LIST_DOMAINS_OTHER          = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_OTHER)
// 	CONNECT_LIST_DOMAINS_MANAGEDSAVE    = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_MANAGEDSAVE)
// 	CONNECT_LIST_DOMAINS_NO_MANAGEDSAVE = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_NO_MANAGEDSAVE)
// 	CONNECT_LIST_DOMAINS_AUTOSTART      = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_AUTOSTART)
// 	CONNECT_LIST_DOMAINS_NO_AUTOSTART   = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_NO_AUTOSTART)
// 	CONNECT_LIST_DOMAINS_HAS_SNAPSHOT   = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_HAS_SNAPSHOT)
// 	CONNECT_LIST_DOMAINS_NO_SNAPSHOT    = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_NO_SNAPSHOT)
// 	CONNECT_LIST_DOMAINS_HAS_CHECKPOINT = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_HAS_CHECKPOINT)
// 	CONNECT_LIST_DOMAINS_NO_CHECKPOINT  = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_NO_CHECKPOINT)
// )
type InstHandle interface{
	LibvirtConnection()
	ReturnDomainList()
}

 

type DomainInfo struct{
		State libvirt.DomainState `json:"state"`
		MaxMem uint64 `json:"maxmem"`
		Memory uint64 `json:"memory"`
		NrVirtCpu uint `json:"nrVirtCpu"`
		CpuTime uint64 `json:"cpuTime"`
		Hwaddr string `json:"hwAddr"`	
}


type PredefinedStructures interface{
	isPredefined()
}
func (D *DomainInfo) isPredefined (){}