package create

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core/services/creation"
	"go.uber.org/zap"
)

type DomainController interface {
	GetDomain(uuid string) (*domCon.Domain, error)
	AddNewDomain(domain *domCon.Domain, uuid string) error
	BootSleepingCPU(domain *domCon.Domain) error
}

type Handler struct {
	creation.VMCreator
	Logger        *zap.Logger
	creation.LibvirtConnect
	DomainControl DomainController
}

func NewHandler(dc DomainController, lc creation.LibvirtConnect, logger *zap.Logger) *Handler {
	return &Handler{
		DomainControl:  dc,
		LibvirtConnect: lc,
		Logger:         logger,
	}
}
