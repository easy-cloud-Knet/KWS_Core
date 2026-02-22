package termination

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"libvirt.org/go/libvirt"
)

func DomainTerminatorFactory(Domain *domCon.Domain) (*DomainTerminator, error) {
	return &DomainTerminator{
		domain: Domain,
	}, nil
}

func (DD *DomainTerminator) TerminateDomain() (*libvirt.Domain, error) {
	dom := DD.domain

	isRunning, err := dom.Domain.IsActive()
	if !isRunning {
		return  nil,virerr.ErrorGen(virerr.DomainShutdownError, fmt.Errorf("error checking domain's aliveness, from libvirt. %w", err))
	}

	if err := dom.Domain.Destroy(); err != nil {
		fmt.Println("error occured while deleting Domain")
		return nil,virerr.ErrorGen(virerr.DomainShutdownError, fmt.Errorf("error shutting down domain, from libvirt. %w, %v", err,DD))
	}

	return dom.Domain, nil
}
