package status

import (
	"fmt"

	virtxml "github.com/easy-cloud-Knet/KWS_Core/pkg/virtXML"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type xmlDom interface {
	GetXMLDesc(flags libvirt.DomainXMLFlags) (string, error)
}

type XMLStatus struct {
	dom xmlDom
}

func (ds *XMLStatus) RetrieveStatus(sources map[SourceType]int, _ *zap.Logger) (map[SourceType]int, error) {
	domainXML, err := virtxml.ConvertExistingDomain(func() (string, error) {
		return ds.dom.GetXMLDesc(0)
	})
	if err != nil {
		return nil, err
	}

	result := make(map[SourceType]int, len(sources))
	for k := range sources {
		switch k {
		case CPU:
			result[CPU] = int(domainXML.VCPU.Value)
		case Memory:
			result[Memory] = int(domainXML.Memory.Value)
		case MaxMemory:
			if domainXML.MaximumMemory != nil {
				result[MaxMemory] = int(domainXML.MaximumMemory.Value)
			} else {
				result[MaxMemory] = int(domainXML.Memory.Value)
			}
		case CPUTime:
			return nil, fmt.Errorf("cpu_time is not available for inactive domains")
		default:
			return nil, fmt.Errorf("unknown source type: %s", string(k))
		}
	}
	return result, nil
}
