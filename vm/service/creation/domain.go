package creation

import (
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

// Configurer generates the XML config bytes for a new VM.
// localConfigurer satisfies this via structural typing.
type Configurer interface {
	GenerateXML(logger *zap.Logger) ([]byte, error)
}

// LibvirtConnect abstracts *libvirt.Connect for domain definition.
// *libvirt.Connect satisfies this via structural typing.
// TODO: DomainDefineXML returns *libvirt.Domain — full mock coverage requires returning Domain interface.
type LibvirtConnect interface {
	DomainDefineXML(xmlConfig string) (*libvirt.Domain, error)
}
