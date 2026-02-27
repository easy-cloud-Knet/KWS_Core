package domListStatus

import (
	"fmt"

	domainStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domain_status"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func (dls *DomainListStatus) UpdateFromDomain(dataDog domainStatus.DataDog, dom *libvirt.Domain, state libvirt.ConnectListAllDomainsFlags, sourceType []domainStatus.SourceType, logger *zap.Logger) error {
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

func (dls *DomainListStatus) GetDomStatus(dom *libvirt.Domain, sourceType []domainStatus.SourceType, logger *zap.Logger) (interface{}, error) {
	state, _, err := dom.GetState()
	if err != nil {
		return nil, fmt.Errorf("failed to get domain state: %w", err)
	}
	// enum 이고 상태값만 봐선 호환될 거 같긴한데
	// 한번 봐야됨, ConnectListAllDomainsFlags <-> DomainState
	dataDog := dls.NewDataDogs(libvirt.ConnectListAllDomainsFlags(state))
	if dataDog == nil {
		return nil, fmt.Errorf("unsupported domain state: %d", state)
	}
	return dataDog.RetrieveStatus(dom, sourceType, logger)

}
