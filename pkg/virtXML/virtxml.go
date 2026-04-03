package virtxml

import libvirtxml "libvirt.org/libvirt-go-xml"

// DomainXMLFlags mirrors libvirt.DomainXMLFlags to avoid a runtime libvirt dependency.
type DomainXMLFlags uint

type DomainXML struct {
	libvirtxml.Domain
}

// New returns a new xml template for domain definition.
// The returned struct is empty and should be filled with necessary information before use.
func New() *libvirtxml.Domain {
	return &libvirtxml.Domain{}
}

// ConvertExistingDomain takes functions as domain.GetXMLDesc or any string based function.
// It should be wrapped with function which returns string and error qualified with libvirt-xml.Domain format.
func ConvertExistingDomain(getXMLDesc func() (string, error)) (*libvirtxml.Domain, error) {
	xmlStr, err := getXMLDesc()
	if err != nil {
		return nil, err
	}

	domain := &libvirtxml.Domain{}
	if err := domain.Unmarshal(xmlStr); err != nil {
		return nil, err
	}

	return domain, nil

}
