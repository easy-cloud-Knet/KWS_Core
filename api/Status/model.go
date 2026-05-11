package status

import (
	svcstatus "github.com/easy-cloud-Knet/KWS_Core/services/status"
)

type DomainStatusRequest struct {
	DataType svcstatus.DomainDataType `json:"dataType"`
	UUID     string                   `json:"UUID"`
}

type HostStatusRequest struct {
	HostDataType svcstatus.HostDataType `json:"host_dataType"`
}

type InstInfoRequest struct {
	InstDataType svcstatus.InstDataType `json:"dataType"`
}

type UUIDListResponse struct {
	UUIDs []string `json:"uuids"`
}

