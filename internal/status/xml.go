package status

import (
	"fmt"

	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
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

func (ds *XMLStatus) RetrieveStatus(sources map[vmtypes.SourceType]int, _ *zap.Logger) (map[vmtypes.SourceType]int, error) {
	domainXML, err := virtxml.ConvertExistingDomain(func() (string, error) {
		return ds.dom.GetXMLDesc(0)
	})
	if err != nil {
		return nil, err
	}
	result := make(map[vmtypes.SourceType]int, len(sources))
	for k := range sources {
		switch k {
		case vmtypes.CPU:
			result[vmtypes.CPU] = int(domainXML.VCPU.Value)
		case vmtypes.Memory:
			result[vmtypes.Memory] = int(domainXML.Memory.Value)
		default:
			return nil, fmt.Errorf("unknown source type: %s", string(k))
		}
	}
	return result, nil
}
