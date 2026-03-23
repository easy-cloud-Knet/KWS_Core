package api

import (
	"github.com/easy-cloud-Knet/KWS_Core/pkg/service/termination"
	"go.uber.org/zap"
)

type ControlHandler struct {
	termination.DomainTermination
	termination.DomainDeleter
	Logger *zap.Logger
}


type DomainControlRequest struct {
	UUID         string                       `json:"UUID"`
	DeletionType termination.DomainDeleteType `json:"DeleteType"`
}
