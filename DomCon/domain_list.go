package domCon

import (
	"fmt"
	"sync"

	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domainList_status"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	instatus "github.com/easy-cloud-Knet/KWS_Core/internal/status"
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
		domainListMutex:  sync.Mutex{},
		DomainList:       make(map[string]*Domain),
		DomainListStatus: &domStatus.DomainListStatus{},
	}
} // 전역적으로 사용되는 도메인 리스트 컨트롤러 생성

func (DC *DomListControl) SetLibvirt(inst *libvirt.Connect) {
	DC.libvirtInst = inst
}

func (DC *DomListControl) AddNewDomain(domain *Domain, uuid string) error {
	DC.domainListMutex.Lock()
	defer DC.domainListMutex.Unlock()

	DC.DomainList[uuid] = domain
	vcpu, err := domain.Domain.GetMaxVcpus()
	if err != nil {
		return virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("error while getting vcpu count during adding new domain: %w", err))
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

func (DC *DomListControl) GetDomain(uuid string) (*Domain, error) {
	DC.domainListMutex.Lock()
	domain, Exist := DC.DomainList[uuid]
	DC.domainListMutex.Unlock()
	if !Exist {
		DomainSeeker := DomSeekUUIDFactory(DC.libvirtInst, uuid)
		dom, err := DomainSeeker.ReturnDomain()
		if err != nil {
			return nil, err
		}
		if err := DC.AddNewDomain(dom, uuid); err != nil {
			return nil, err
		}
		return dom, nil
	}

	return domain, nil
}

func (DC *DomListControl) SleepDomain(domain *Domain, logger *zap.Logger) error {
	sources := map[instatus.SourceType]int{instatus.CPU: 0}
	stat, err := DC.DomainListStatus.GetDomStatus(domain.Domain, sources, logger)
	if err != nil {
		return err
	}
	DC.DomainListStatus.AddSleepingCPU(stat[instatus.CPU])
	return nil
}

func (DC *DomListControl) RemoveDomain(domain *Domain, uuid string, logger *zap.Logger) error {
	sources := map[instatus.SourceType]int{instatus.CPU: 0}
	stat, err := DC.DomainListStatus.GetDomStatus(domain.Domain, sources, logger)
	if err != nil {
		return err
	}
	vcpu := stat[instatus.CPU]

	DC.domainListMutex.Lock()
	delete(DC.DomainList, uuid)
	domain.Domain.Free()
	DC.domainListMutex.Unlock()
	DC.DomainListStatus.TakeAllocatedCPU(vcpu)
	return nil
}

func (DC *DomListControl) retrieveDomainsByState(state libvirt.ConnectListAllDomainsFlags, logger *zap.Logger) error {
	domains, err := DC.libvirtInst.ListAllDomains(state)
	if err != nil {
		logger.Fatal("Failed to retrieve domains", zap.Error(err))
		return err
	}

	isActive := state == libvirt.CONNECT_LIST_DOMAINS_ACTIVE
	wg := &sync.WaitGroup{}
	for _, dom := range domains {
		uuid, err := dom.GetUUIDString()
		if err != nil {
			logger.Sugar().Error("Failed to get UUID for domain", err)
			continue
		}
		NewDom := &Domain{
			Domain:      &dom,
			domainMutex: sync.Mutex{},
		}
		DC.AddExistingDomain(NewDom, uuid)

		wg.Add(1)
		go func(targetDom libvirt.Domain, targetUUID string) {
			defer wg.Done()
			dataDog := instatus.New(&targetDom, isActive)
			sources := map[instatus.SourceType]int{instatus.CPU: 0}
			if err := DC.DomainListStatus.UpdateFromDomain(dataDog, isActive, sources, logger); err != nil {
				logger.Sugar().Errorf("Failed to retrieve status for domain UUID=%s: %v", targetUUID, err)
			}
		}(dom, uuid)
		logger.Sugar().Infof("Added domain: UUID=%s", uuid)
	}
	wg.Wait()

	logger.Sugar().Infof("Total %d domains added (state: %d)", len(domains), state)
	return nil
}

func (DC *DomListControl) RetrieveAllDomain(logger *zap.Logger) error {
	logger.Info("Retrieving all domains from libvirt...")

	if err := DC.retrieveDomainsByState(libvirt.CONNECT_LIST_DOMAINS_ACTIVE, logger); err != nil {
		return err
	}

	if err := DC.retrieveDomainsByState(libvirt.CONNECT_LIST_DOMAINS_INACTIVE, logger); err != nil {
		return err
	}

	logger.Info("retreiving intital vm", zap.Int("number", len(DC.DomainList)))
	return nil
}

func (DC *DomListControl) GetDomainListStatus() *domStatus.DomainListStatus {
	return DC.DomainListStatus
}

////////////////////////////////////////////////

func (DC *DomListControl) BootSleepingCPU(domain *Domain) error {
	vcpu, err := domain.Domain.GetMaxVcpus()
	if err != nil {
		return err
	}
	DC.DomainListStatus.TakeSleepingCPU(int(vcpu))
	return nil
}

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
