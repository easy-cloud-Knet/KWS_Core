package conn

import (
	"fmt"
	"sync"

	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func NewDomainInstance(Dom *libvirt.Domain) *Domain {
	return &Domain{
		domainMutex: sync.Mutex{},
		Domain:      Dom,
	}
}

func DomListConGen() *DomListControl {
	return &DomListControl{
		domainMutex: sync.Mutex{},
		DomainList:  make(map[string]*Domain),
	}
}

func (DC *DomListControl) AddNewDomain(domain *Domain, uuid string) {
	DC.domainMutex.Lock()
	defer DC.domainMutex.Unlock()

	DC.DomainList[uuid] = domain
}

func (DC *DomListControl) GetDomain(uuid string, LibvirtInst *libvirt.Connect) (*Domain, error) {
	DC.domainMutex.Lock()
	domain, Exist := DC.DomainList[uuid]
	DC.domainMutex.Unlock()

	if !Exist {
		DomainSeeker := DomSeekUUIDFactory(LibvirtInst, uuid)
		domList, err := DomainSeeker.ReturnDomain()
		if err != nil {
			return nil, err
		}
		DC.AddNewDomain(domList, uuid)
		return domList, nil
	}

	return domain, nil
}

func (DC *DomListControl) DeleteDomain(uuid string, LibvirtInst *libvirt.Connect) error {
	DC.domainMutex.Lock()
	domain, Exist := DC.DomainList[uuid]
	DC.domainMutex.Unlock()

	if !Exist {
		DomainSeeker := DomSeekUUIDFactory(LibvirtInst, uuid)
		dom, err := DomainSeeker.ReturnDomain()
		if err != nil {
			return virerr.ErrorGen(virerr.NoSuchDomain, fmt.Errorf("domain trying to delete already empty, uuid of %s, %w", uuid, err))
		}
		dom.Domain.Free()
		return nil
	}

	domain.Domain.Free()

	DC.domainMutex.Lock()
	delete(DC.DomainList, uuid)
	DC.domainMutex.Unlock()

	return nil
}

func (DC *DomListControl) retrieveDomainsByState(LibvirtInst *libvirt.Connect, state libvirt.ConnectListAllDomainsFlags, logger *zap.Logger) error {
	domains, err := LibvirtInst.ListAllDomains(state)
	if err != nil {
		logger.Fatal("Failed to retrieve domains",zap.Error(err))
		return err
	}

	DC.domainMutex.Lock()
	defer DC.domainMutex.Unlock()

	for _, dom := range domains {
		uuid, err := dom.GetUUIDString()
		if err != nil {
			logger.Sugar().Error("Failed to get UUID for domain", err)
			continue
		}

		DC.DomainList[uuid] = &Domain{
			Domain:      &dom,
			domainMutex: sync.Mutex{},
		}
		// logger.Infof("Added domain: UUID=%s", uuid)
		logger.Sugar().Infof("Added domain: UUID=%s", uuid)
	}

	logger.Sugar().Infof("Total %d domains added (state: %d)", len(domains), state)
	return nil
}

func (DC *DomListControl) RetrieveAllDomain(LibvirtInst *libvirt.Connect, logger *zap.Logger) error {
	logger.Info("Retrieving all domains from libvirt...")

	if err := DC.retrieveDomainsByState(LibvirtInst, libvirt.CONNECT_LIST_DOMAINS_ACTIVE, logger); err != nil {
		return err
	}

	if err := DC.retrieveDomainsByState(LibvirtInst, libvirt.CONNECT_LIST_DOMAINS_INACTIVE, logger); err != nil {
		return err
	}

	logger.Info("retreiving intital vm", zap.Int("number", len(DC.DomainList)))
	return nil
}
