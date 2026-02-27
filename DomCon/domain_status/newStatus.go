package domainStatus

import (
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func (ds *XMLStatus) RetrieveStatus(dom *libvirt.Domain, sources []SourceType, logger *zap.Logger) (interface{}, error) {
	domcnf, err := XMLUnparse(dom)
	if err != nil {
		logger.Error("failed to unparse domain XML", zap.Error(err))
		return nil, err
	}
	mapSource := make(map[SourceType]int)
	for _, source := range sources {
		switch source {
		case CPU:
			mapSource[CPU] = int(domcnf.VCPU.Value)
		case Memory:
			mapSource[Memory] = int(domcnf.Memory.Value)
		default:
			logger.Warn("unknown source type", zap.String("source", string(source)))
		}
	}
	return mapSource, nil

}

func (ls *LibvirtStatus) RetrieveStatus(dom *libvirt.Domain, sources []SourceType, logger *zap.Logger) (interface{}, error) {

	mapSource := make(map[SourceType]int)
	for _, source := range sources {
		switch source {
		case CPU:
			cpu, err := ls.RetrieveCPU(dom, *logger)
			if err != nil {
				logger.Error("failed to retrieve CPU count", zap.Error(err))
				return nil, err
			}
			mapSource[CPU] = cpu
		default:
			logger.Warn("unknown source type", zap.String("source", string(source)))
		}
	}
	return mapSource, nil

}
