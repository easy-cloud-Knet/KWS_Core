package status

import (
	"fmt"

	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	"go.uber.org/zap"
)

type libvirtDom interface {
	GetMaxVcpus() (uint, error)
}

type LibvirtStatus struct {
	dom libvirtDom
}

func (ls *LibvirtStatus) RetrieveStatus(sources map[vmtypes.SourceType]int, _ *zap.Logger) (map[vmtypes.SourceType]int, error) {
	for k := range sources {
		switch k {
		case vmtypes.CPU:
			cpu, err := ls.dom.GetMaxVcpus()
			if err != nil {
				return nil, err
			}
			sources[vmtypes.CPU] = int(cpu)
		default:
			return nil, fmt.Errorf("unknown source type: %s", string(k))
		}
	}
	return sources, nil
}
