package snapshot

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"go.uber.org/zap"
)

type DomainController interface {
	GetDomain(uuid string) (*domCon.Domain, error)
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
