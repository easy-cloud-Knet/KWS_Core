package control

import "github.com/easy-cloud-Knet/KWS_Core/services/termination"

type DomainControlRequest struct {
	UUID         string                       `json:"UUID"`
	DeletionType termination.DomainDeleteType `json:"DeleteType"`
}
