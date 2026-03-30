package domCon

import (
	"errors"
	"sync"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	pkguuid "github.com/easy-cloud-Knet/KWS_Core/pkg/UUID"
	"libvirt.org/go/libvirt"
)

func DomSeekUUIDFactory(LibInstance *libvirt.Connect, UUID string) *DomainSeekingByUUID {
	return &DomainSeekingByUUID{
		LibvirtInst: LibInstance,
		UUID:        UUID,
	}
}

func (DSU *DomainSeekingByUUID) ReturnDomain() (*Domain, error) {
	parsedUUID, err := pkguuid.ValidateAndReturnUUID(DSU.UUID)
	if err != nil {
		return nil, virerr.ErrorGen(virerr.InvalidUUID, err)
	}
	domain, err := DSU.LibvirtInst.LookupDomainByUUID((*parsedUUID)[:])
	if err != nil {
		var libErr libvirt.Error
		if errors.As(err, &libErr) && libErr.Code == libvirt.ERR_NO_DOMAIN {
			return nil, virerr.ErrorGen(virerr.NoSuchDomain, err)
		}
		return nil, virerr.ErrorGen(virerr.DomainSearchError, err)
	} else if domain == nil {
		return nil, virerr.ErrorGen(virerr.NoSuchDomain, err)
	}

	return &Domain{
		Domain:      domain,
		domainMutex: sync.Mutex{},
	}, nil
}
