package conn

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core.git/api/conn/DomCon"
	"libvirt.org/go/libvirt"
)


type DomainController struct{
	DomainList *domCon.DomListControl
	Operator Operator
}

type Operator interface{
	Operation() (*libvirt.Domain,error)
}


