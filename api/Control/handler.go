package control

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"go.uber.org/zap"
)

type DomainController interface {
	GetDomain(uuid string) (*domCon.Domain, error)
	SleepDomain(domain *domCon.Domain, logger *zap.Logger) error
	RemoveDomain(domain *domCon.Domain, uuid string, logger *zap.Logger) error
}

type Handler struct {
	DomainControl DomainController
	Logger        *zap.Logger
}

func NewHandler(dc DomainController, logger *zap.Logger) *Handler {
	return &Handler{
		DomainControl: dc,
		Logger:        logger,
	}
}
