package api

import (
	"github.com/easy-cloud-Knet/KWS_Core/vm/service/status"
	"libvirt.org/go/libvirt"
)

//// income api Structures

type ShutDownDomain struct {
	UUID string `json:"UUID"`
}
type StartDomain struct {
	UUID string `json:"UUID"`
}

type ReturnDomainFromUUID struct {
	DataType status.DomainDataType `json:"dataType"`
	UUID     string                `json:"UUID"`
}

// host
type ReturnHostFromStatus struct {
	HostDataType status.HostDataType `json:"host_dataType"`
}

type ReturnInstAllData struct {
	InstDataType status.InstDataType `json:"dataType"`
}

// //////////////////////
type UUIDListResponse struct {
	UUIDs []string `json:"uuids"`
}

type DomainState_init struct {
	DomainState libvirt.DomainState `json:"currentState"`
	UUID        string              `json:"UUID"`
}
