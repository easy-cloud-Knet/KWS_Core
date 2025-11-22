package domCon

import (
	"fmt"
	"sync"

	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domain_status"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
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
		domainListMutex: sync.Mutex{},
		DomainList:      make(map[string]*Domain),
		DomainListStatus: &domStatus.DomainListStatus{},
	}
}

func (DC *DomListControl) AddNewDomain(domain *Domain, uuid string) error {
	DC.domainListMutex.Lock()
	defer DC.domainListMutex.Unlock()

	DC.DomainList[uuid] = domain
	vcpu, err :=domain.Domain.GetMaxVcpus()
	if err != nil {
		Err:=fmt.Errorf("%v error while getting vcpu count during adding new domain",err)
		return Err
	}
	DC.DomainListStatus.AddAllocatedCPU(int(vcpu))
	return nil
}

func (DC *DomListControl) AddExistingDomain(domain *Domain, uuid string) {
	DC.domainListMutex.Lock()
	defer DC.domainListMutex.Unlock()

	DC.DomainList[uuid] = domain
}
// Exstring Domain only called from initial booting, and adding specs is not its role

func (DC *DomListControl) GetDomain(uuid string, LibvirtInst *libvirt.Connect) (*Domain, error) {
	fmt.Println(DC)
	DC.domainListMutex.Lock()
	domain, Exist := DC.DomainList[uuid]
	DC.domainListMutex.Unlock()
	if !Exist {
		DomainSeeker := DomSeekUUIDFactory(LibvirtInst, uuid)
		dom, err := DomainSeeker.ReturnDomain()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		fmt.Println(dom)
		DC.AddNewDomain(dom, uuid)
		return dom, nil
	}
	fmt.Println(domain)	

	return domain, nil
}

func (DC *DomListControl) DeleteDomain(Domain *libvirt.Domain, uuid string, vcpu int) error {
	DC.domainListMutex.Lock()
	delete(DC.DomainList, uuid)
	Domain.Free()
	DC.domainListMutex.Unlock()
	DC.DomainListStatus.TakeAllocatedCPU(vcpu)
	return nil
}

func (DC *DomListControl) FindAndDeleteDomain(LibvirtInst *libvirt.Connect, uuid string) error {
	//아직 활용처가 없어서, vcpu 삭제를 추가하지 않았음/
	// DeleteDomain 함수의 TakeAllocatedCPU 호출을 참고..
	DC.domainListMutex.Lock()
	domain, Exist := DC.DomainList[uuid]
	DC.domainListMutex.Unlock()

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

	DC.domainListMutex.Lock()
	delete(DC.DomainList, uuid)
	DC.domainListMutex.Unlock()

	return nil
}

func (DC *DomListControl) retrieveDomainsByState(LibvirtInst *libvirt.Connect, state libvirt.ConnectListAllDomainsFlags, logger *zap.Logger) error {
	domains, err := LibvirtInst.ListAllDomains(state)
	if err != nil {
		logger.Fatal("Failed to retrieve domains", zap.Error(err))
		return err
	}


	dataDog := domStatus.NewDataDog(state)
	wg:= &sync.WaitGroup{}
	for _, dom := range domains {
		uuid, err := dom.GetUUIDString()
		if err != nil {
			logger.Sugar().Error("Failed to get UUID for domain", err)
			continue
		}
		NewDom:= &Domain{
			Domain:      &dom,
			domainMutex: sync.Mutex{},
		}
		DC.AddExistingDomain(NewDom,uuid)
		
		wg.Add(1)
		go func(targetDom libvirt.Domain) { 
		defer wg.Done()
		retrieveFunc := dataDog.Retreive(&targetDom, DC.DomainListStatus, *logger)
		if retrieveFunc != nil {
			logger.Sugar().Errorf("Failed to retrieve status for domain UUID=%s: %v", uuid, retrieveFunc)
		}

		}(dom)
		logger.Sugar().Infof("Added domain: UUID=%s", uuid)
	}
	wg.Wait()
	fmt.Printf("%+v", *DC.DomainListStatus)
	
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

////////////////////////////////////////////////

func (DC *DomListControl) GetAllUUIDs() []string {
	DC.domainListMutex.Lock()
	defer DC.domainListMutex.Unlock()

	uuids := make([]string, 0, len(DC.DomainList))
	for uuid := range DC.DomainList {
		uuids = append(uuids, uuid)
	}
	return uuids
}

//////////////////////////////////////////////// 