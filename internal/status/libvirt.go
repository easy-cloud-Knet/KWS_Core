package status

import (
	"fmt"

	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type libvirtDom interface {
	GetMaxVcpus() (uint, error)
	GetInfo() (*libvirt.DomainInfo, error)
	GetMaxMemory() (uint64, error)
}

type LibvirtStatus struct {
	dom libvirtDom
}

func (ls *LibvirtStatus) RetrieveStatus(sources map[SourceType]int, _ *zap.Logger) (map[SourceType]int, error) {
	var info *libvirt.DomainInfo
	getInfo := func() (*libvirt.DomainInfo, error) {
		if info != nil {
			return info, nil
		}
		var err error
		info, err = ls.dom.GetInfo()
		return info, err
	}

	result := make(map[SourceType]int, len(sources))
	for k := range sources {
		switch k {
		case CPU:
			cpu, err := ls.dom.GetMaxVcpus()
			if err != nil {
				return nil, err
			}
			result[CPU] = int(cpu)
		case Memory:
			i, err := getInfo()
			if err != nil {
				return nil, err
			}
			result[Memory] = int(i.Memory)
		case MaxMemory:
			mem, err := ls.dom.GetMaxMemory()
			if err != nil {
				return nil, err
			}
			result[MaxMemory] = int(mem)
		case CPUTime:
			i, err := getInfo()
			if err != nil {
				return nil, err
			}
			result[CPUTime] = int(i.CpuTime)
		default:
			return nil, fmt.Errorf("unknown source type: %s", string(k))
		}
	}
	return result, nil
}
