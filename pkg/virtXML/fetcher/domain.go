package fetcher

import (
	"fmt"

	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	virtxml "github.com/easy-cloud-Knet/KWS_Core/pkg/virtXML"
	"libvirt.org/go/libvirt"
	libvirtxml "libvirt.org/libvirt-go-xml"
)

type XMLFetcher struct{}

type Domain interface {
	GetXMLDesc(flags libvirt.DomainXMLFlags) (string, error)
}

func NewXMLFetcher() *XMLFetcher {
	return &XMLFetcher{}
}

func (xf *XMLFetcher) Fetch(domain Domain, sources map[vmtypes.SourceType]int) (map[vmtypes.SourceType]int, error) {
	domainXML, err := xf.parse(domain)
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

func (xf *XMLFetcher) parse(domain Domain) (*libvirtxml.Domain, error) {
	return virtxml.ConvertExistingDomain(func() (string, error) {
		return domain.GetXMLDesc(0)
	})
}
