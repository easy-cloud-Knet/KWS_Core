package status

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domainList_status"
	svcstatus "github.com/easy-cloud-Knet/KWS_Core/services/status"
	"go.uber.org/zap"
)

type DomainController interface {
	GetDomain(uuid string) (*domCon.Domain, error)
	GetAllUUIDs() []string
	GetDomainListStatus() *domStatus.DomainListStatus
}

type Handler struct {
	LibvirtConn   svcstatus.Connect
	DomainControl DomainController
	Logger        *zap.Logger
}

func NewHandler(lv svcstatus.Connect, dc DomainController, logger *zap.Logger) *Handler {
	return &Handler{
		LibvirtConn:   lv,
		DomainControl: dc,
		Logger:        logger,
	}
}
