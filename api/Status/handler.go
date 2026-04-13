package status

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domainList_status"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type DomainController interface {
	GetDomain(uuid string) (*domCon.Domain, error)
	GetAllUUIDs() []string
	GetDomainListStatus() *domStatus.DomainListStatus
}

type Handler struct {
	LibvirtInst   *libvirt.Connect
	DomainControl DomainController
	Logger        *zap.Logger
}

func NewHandler(lv *libvirt.Connect, dc DomainController, logger *zap.Logger) *Handler {
	return &Handler{
		LibvirtInst:   lv,
		DomainControl: dc,
		Logger:        logger,
	}
}
