package domCon

import (
	"sync"

	"libvirt.org/go/libvirt"
)

// libvirt 를 통해 도메인을 찾는 행위를 최소한 하기위해 관리하는 리스트(찾을때 걸리는 시간, 사용 후 도메인을 해제하는 과정 최소화)
// libvirt 내에서 domain은 *libvirt.domain으로 관리 됨
// DomainList 에서 uuid형태로 각각의 도메인을 관리


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
