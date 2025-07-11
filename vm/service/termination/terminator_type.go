package termination

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"libvirt.org/go/libvirt"
)

type DomainDeleteType uint

const (
	HardDelete DomainDeleteType = iota
	SoftDelete
)

type DomainDeletion interface{
	DeleteDomain() (*libvirt.Domain,error)
}

type DomainTermination interface{
	TerminateDomain() (*libvirt.Domain,error)
}

type DomainTerminator struct {
	domain *domCon.Domain
}
type DomainDeleter struct {
	uuid string
	domain        *domCon.Domain
	DeletionType        DomainDeleteType
}
