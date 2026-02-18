package api

import "github.com/easy-cloud-Knet/KWS_Core/vm/service/termination"

type DeleteDomain struct {
	UUID         string                       `json:"UUID"`
	DeletionType termination.DomainDeleteType `json:"DeleteType"`
}
