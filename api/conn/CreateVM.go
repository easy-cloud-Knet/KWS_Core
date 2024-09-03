package conn

import (
	"log"
	"net/http"

	"libvirt.org/go/libvirt"
)

type Domain struct{
	Domain *libvirt.Domain
}

type  BasicDomainControl interface{
	createDomain()
}

func (d *Domain)CreateVM(w *http.ResponseWriter, r * http.Request){
	log.Fatal(d.Domain.Create())
}
