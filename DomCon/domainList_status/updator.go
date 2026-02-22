package domListStatus

import (
	"fmt"

	domainStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domain_status"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func (dls *DomainListStatus) UpdateFromDomain(dataDog domainStatus.DataDog, dom *libvirt.Domain, state libvirt.ConnectListAllDomainsFlags, sourceType []domainStatus.SourceType, logger zap.Logger) error {
	result, err := dataDog.RetrieveStatus(dom, sourceType, logger)
	if err != nil {
		return err
	}
	statusMap, ok := result.(map[domainStatus.SourceType]int)
	if !ok {
		return fmt.Errorf("unexpected status type from domain retrieval")
	}
	if cpu, ok := statusMap[domainStatus.CPU]; ok {
		dls.AddAllocatedCPU(cpu)
		if state == libvirt.CONNECT_LIST_DOMAINS_INACTIVE {
			dls.AddSleepingCPU(cpu)
		}
	}
	return nil
}

func (dls *DomainListStatus) NewDataDogs(state libvirt.ConnectListAllDomainsFlags) domainStatus.DataDog {
	switch state {
	case libvirt.CONNECT_LIST_DOMAINS_ACTIVE:
		return &domainStatus.LibvirtStatus{}
	case libvirt.CONNECT_LIST_DOMAINS_INACTIVE:
		return &domainStatus.XMLStatus{}
	default:
		return nil
	}
}
