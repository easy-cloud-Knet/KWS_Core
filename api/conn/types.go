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