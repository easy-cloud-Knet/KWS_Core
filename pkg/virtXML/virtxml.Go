package virtxml

import libvirtxml "libvirt.org/libvirt-go-xml"

type DomainXML struct {
	libvirtxml.Domain
}

// New returns a new xml template for domain definition.
// The returned struct is empty and should be filled with necessary information before use.
func New() *libvirtxml.Domain {
	return &libvirtxml.Domain{}
}

// ConvertExistingDomain takes libvirt.XMLDesc() into libvirtxml.Domain struct.
func ConvertExistingDomain(Dflag DomainXMLFlags, getXMLDesc func(flags DomainXMLFlags) (string, error)) (*libvirtxml.Domain, error) {
	xmlStr, err := getXMLDesc(DomainXMLFlags(Dflag))
	if err != nil {
		return nil, err
	}

	domain := &libvirtxml.Domain{}
	if err := domain.Unmarshal(xmlStr); err != nil {
		return nil, err
	}

	return domain, nil

}
