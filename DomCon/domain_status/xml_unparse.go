package domStatus

import (
	"fmt"

	"libvirt.org/go/libvirt"
	libvirtxml "libvirt.org/libvirt-go-xml"
)




func XMLUnparse(domain * libvirt.Domain) error {
	
	domainXML, err := domain.GetXMLDesc(0)
	if err != nil {
		return err
	}
	domcnf := &libvirtxml.Domain{}

	err= domcnf.Unmarshal(domainXML)
	if err != nil {
		return err
	}
	fmt.Println(domcnf)
	return nil
}

