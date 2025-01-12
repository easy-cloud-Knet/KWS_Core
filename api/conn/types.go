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

type DataTypeHandler interface{
	GetInfo(*Domain) error
	// Generator(DomainDataType) err
}
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
	Hwaddr string `json:"hwAddr"`	
}

type InstHandler struct{
	LibvirtInst *libvirt.Connect
}

type DomainDetail struct{
	DomainSeeker DomainSeeker
	DataHandle []DataTypeHandler 
}

type DomainSeekinggByUUID struct{
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
	returnDomain()([]*Domain,error)
}

type ReturnDomainFromStatus struct{ 
	DataType DomainDataType `json:"type"`
	Status libvirt.ConnectListAllDomainsFlags  `json:"Status Flag"`
}

type ReturnDomainFromUUID struct{ 
	DataType DomainDataType `json:"type"`
	UUID string  `json:"Status UUID"`
}

type InstHandle interface{
	LibvirtConnection()
}




 