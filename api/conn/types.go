package conn

import (
	"sync"

	"libvirt.org/go/libvirt"
)




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






type DomainInfo struct{
		State libvirt.DomainState `json:"state"`
		MaxMem uint64 `json:"maxmem"`
		Memory uint64 `json:"memory"`
		NrVirtCpu uint `json:"nrVirtCpu"`
		CpuTime uint64 `json:"cpuTime"`
		Hwaddr string `json:"hwAddr"`	
		UUID string `json:uuid`
}
