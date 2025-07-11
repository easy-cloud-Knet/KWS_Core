package domCon

import (
	"sync"

	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"github.com/google/uuid"
	"libvirt.org/go/libvirt"
)


func DomSeekUUIDFactory(LibInstance *libvirt.Connect,UUID string)*DomainSeekingByUUID{
	return &DomainSeekingByUUID{ 
		LibvirtInst: LibInstance,
		UUID:        UUID,
	}
}


func ReturnUUID(UUID string) (*uuid.UUID, error) {
	uuidParsed, err := uuid.Parse(UUID)
	if err != nil {
		return nil, err
	}
	return &uuidParsed, nil
}


func (DSU *DomainSeekingByUUID) ReturnDomain() (*Domain, error) {
	parsedUUID, err := uuid.Parse(DSU.UUID)
	if err != nil {
		return nil,virerr.ErrorGen(virerr.InvalidUUID, err)
	}
	domain, err := DSU.LibvirtInst.LookupDomainByUUID(parsedUUID[:])
	if err != nil {
		return nil,virerr.ErrorGen(virerr.DomainSearchError, err)
	}else if domain==nil {
		return nil,virerr.ErrorGen(virerr.NoSuchDomain, err)
	}

	
	return &Domain{
		Domain:      domain,
		domainMutex: sync.Mutex{},
	}, nil
}

