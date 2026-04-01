package domListStatus

import (
	"fmt"

	instatus "github.com/easy-cloud-Knet/KWS_Core/internal/status"
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func (dls *DomainListStatus) UpdateFromDomain(dataDog instatus.StatusRetriever, isActive bool, sources map[vmtypes.SourceType]int, logger *zap.Logger) error {
	statusMap, err := dataDog.RetrieveStatus(sources, logger)
	if err != nil {
		return err
	}
	if cpu, ok := statusMap[vmtypes.CPU]; ok {
		dls.AddAllocatedCPU(cpu)
		if !isActive {
			dls.AddSleepingCPU(cpu)
		}
	}
	return nil
}

func (dls *DomainListStatus) GetDomStatus(dom *libvirt.Domain, sources map[vmtypes.SourceType]int, logger *zap.Logger) (map[vmtypes.SourceType]int, error) {
	isActive, err := dom.IsActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get domain state: %w", err)
	}
	return instatus.New(dom, isActive).RetrieveStatus(sources, logger)
}
