package domCon

import (
	"sync"

	"libvirt.org/go/libvirt"
)


type DomListControl struct {
	DomainList map[string]*Domain
	domainListMutex sync.Mutex 
}

type Domain struct {
	Domain      *libvirt.Domain 
	domainMutex sync.Mutex 
}



type DomainSeekingByUUID struct {
	LibvirtInst *libvirt.Connect
	UUID        string
}

type DomainSeeker interface {
	ReturnDomain() (*Domain, error)
}
