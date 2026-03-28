package status

import (
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	"go.uber.org/zap"
)

type LibvirtStatus struct{}

func (ls *LibvirtStatus) RetrieveStatus(dom Domain, sources map[vmtypes.SourceType]int, logger *zap.Logger) (map[vmtypes.SourceType]int, error) {
	for k := range sources {
		switch k {
		case vmtypes.CPU:
			cpu, err := dom.GetMaxVcpus()
			if err != nil {
				return nil, err
			}
			sources[vmtypes.CPU] = int(cpu)
		default:
			logger.Warn("unknown source type", zap.String("source", string(k)))
		}
	}
	return sources, nil
}
