package termination

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core.git/api/conn/DomCon"
)

type DomainDeleteType uint

const (
	HardDelete DomainDeleteType = iota
	SoftDelete
)

type DomainTerminator struct {
	domain *domCon.Domain
}
type DomainDeleter struct {
	uuid string
	domain        *domCon.Domain
	DeletionType        DomainDeleteType
}
