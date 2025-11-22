package domStatus

import (
	"fmt"

	"libvirt.org/go/libvirt"
	libvirtxml "libvirt.org/libvirt-go-xml"
)

// 꺼져있는 도메인의 xml 을 파싱하여 도메인 상태를 업데이트
func XMLUnparse(domain * libvirt.Domain) (*libvirtxml.Domain, error) {
	
	domainXML, err := domain.GetXMLDesc(0)
	if err != nil {
		return nil,fmt.Errorf("%xerror occured while calling xml specification", err)
	}
	domcnf := &libvirtxml.Domain{}

	err= domcnf.Unmarshal(domainXML)
	if err != nil {
		return nil, fmt.Errorf("%x error occured while unmarshalling xml, check for crushed format", err)
	}

	return domcnf, nil
}

