package conn

import (
	"log"

	"libvirt.org/go/libvirt"
)

type Domain struct{
	Domain *libvirt.Domain
}

type  BasicDomainControl interface{
	createDomain()
}

func (d *Domain)CreateVM(){
	log.Fatal(d.Domain.Create())
}
