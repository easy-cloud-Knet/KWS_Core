package api

import (
	"github.com/easy-cloud-Knet/KWS_Core/vm/service/status"
	"libvirt.org/go/libvirt"
)

type DomainBootRequest struct {
	UUID string `json:"UUID"`
}

type DomainStatusRequest struct {
	DataType status.DomainDataType `json:"dataType"`
	UUID     string                `json:"UUID"`
}

type HostStatusRequest struct {
	HostDataType status.HostDataType `json:"host_dataType"`
}

type InstInfoRequest struct {
	InstDataType status.InstDataType `json:"dataType"`
}

type UUIDListResponse struct {
	UUIDs []string `json:"uuids"`
}

type DomainStateResponse struct {
	DomainState libvirt.DomainState `json:"currentState"`
	UUID        string              `json:"UUID"`
}
