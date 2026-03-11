package status

import "libvirt.org/go/libvirt"

// Domain is a minimal interface over *libvirt.Domain for status service.
// *libvirt.Domain satisfies this via structural typing.
type Domain interface {
	GetInfo() (*libvirt.DomainInfo, error)
	GetState() (libvirt.DomainState, int, error)
	GetUUID() ([]byte, error)
	GetGuestInfo(types libvirt.DomainGuestInfoTypes, flags uint32) (*libvirt.DomainGuestInfo, error)
}

// Connect is a minimal interface over *libvirt.Connect for listing all domains.
type Connect interface {
	ListAllDomains(flags libvirt.ConnectListAllDomainsFlags) ([]libvirt.Domain, error)
}
