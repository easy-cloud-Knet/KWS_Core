package domainStatus

import (
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func (ls *LibvirtStatus) RetrieveCPU(dom *libvirt.Domain, logger zap.Logger) (int, error) {
	cpuCount, err := dom.GetMaxVcpus()
	if err != nil {
		logger.Error("failed to get live vcpu count", zap.Error(err))
		return 0, err
	}

	return int(cpuCount), nil

}
