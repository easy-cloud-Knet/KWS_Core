package conn

import (
	"sync"

	"libvirt.org/go/libvirt"
)
type DomainDataType uint

const (
	PowerStaus DomainDataType =iota
	BasicInfo 	
)

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

type DataType struct{
	DataType DomainDataType `json:"type"`
	LibvirtInst *libvirt.Connect
}
type DataTypeHandle interface{
	Setter()
	Getter() (DomainDataType,*libvirt.Connect)
}
type InstHandler struct{
	LibvirtInst *libvirt.Connect
}

type DomainSortingByUUID struct{
	TypeBase DataType 
	UUID string `json:"UUID"`
}

type DomainSortingByStatus struct{
	TypeBase DataType 
	Status libvirt.ConnectListAllDomainsFlags `json:"Status Flag"`
}

type DomainSeeker interface{
	returnStatus() ([]*libvirt.Domain, error)
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
		UUID string `json:"uuid"`
}
