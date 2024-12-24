package conn

import (
	"sync"

	"libvirt.org/go/libvirt"
)


type IP []byte


type DomainList struct{
	RequestType string `json:"requestType"` 
	// libvirt.ConnectListAllDomainsFlags
}


type Domain struct{
	Domain *libvirt.Domain
}

type  BasicDomainControl interface{
	createDomain()
}


type InstHandler struct{
	LibvirtInst *libvirt.Connect
	InstMu sync.Mutex
}

type InstHandle interface{
	LibvirtConnection()
	ActiveDomain()
	ReturnDomainList()
}

type Create_VM_Method int8
const (
	CREATE_WITH_XML Create_VM_Method = iota+1
	type1
	type2
	type3
)


type VM_Init_Info struct{
	RAM int `json:"RAM"`
	CPU int `json:"CPU"`
	IP string `json:"IP"`
	Method Create_VM_Method `json:"METHOD"`
}



type DomainInfo struct{
		State libvirt.DomainState `json:"state"`
		MaxMem uint64 `json:"maxmem"`
		Memory uint64 `json:"memory"`
		NrVirtCpu uint `json:"nrVirtCpu"`
		CpuTime uint64 `json:"cpuTime"`
		Hwaddr string `json:"hwAddr"`	
		UUID string `json:uuid`
}
