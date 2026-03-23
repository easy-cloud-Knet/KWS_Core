package termination

import (
	"fmt"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func DomainTerminatorFactory(d Domain) DomainTermination {
	return &DomainTerminator{
		domain: d,
	}
}

func (DD *DomainTerminator) TerminateDomain() error {
	isRunning, err := DD.domain.IsActive()
	if !isRunning {
		return virerr.ErrorGen(virerr.DomainShutdownError, fmt.Errorf("error checking domain's aliveness, from libvirt. %w", err))
	}

	if err := DD.domain.Destroy(); err != nil {
		return virerr.ErrorGen(virerr.DomainShutdownError, fmt.Errorf("error shutting down domain, from libvirt. %w", err))
	}

	return nil
}
