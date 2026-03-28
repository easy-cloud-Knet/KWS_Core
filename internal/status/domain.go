package status

import "libvirt.org/go/libvirt"

type Domain interface {
	GetMaxVcpus() (uint, error)
	GetXMLDesc(flags libvirt.DomainXMLFlags) (string, error)
}
