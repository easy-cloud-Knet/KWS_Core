package api

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core/pkg/service/termination"
	"go.uber.org/zap"
)

type ControlHandler struct {
	termination.DomainTermination
	termination.DomainDeleter
	Logger *zap.Logger
}

type terminantion_listor interface {
}

func instance_converter(DomainListControl *domCon.DomListControl) terminantion_listor {
	return nil
}

type DomainControlRequest struct {
	UUID         string                       `json:"UUID"`
	DeletionType termination.DomainDeleteType `json:"DeleteType"`
}
