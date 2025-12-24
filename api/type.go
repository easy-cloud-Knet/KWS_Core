package api

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	snapmgr "github.com/easy-cloud-Knet/KWS_Core/vm/service/snapshot"
	"github.com/easy-cloud-Knet/KWS_Core/vm/service/status"
	"github.com/easy-cloud-Knet/KWS_Core/vm/service/termination"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type InstHandler struct {
	LibvirtInst     *libvirt.Connect
	DomainControl   *domCon.DomListControl
	Logger          *zap.Logger
	SnapshotManager snapmgr.SnapshotManager
}

// InstHandler 는

type InstHandle interface {
	LibvirtConnection()
}

//// income api Structures

type DeleteDomain struct {
	UUID         string                       `json:"UUID"`
	DeletionType termination.DomainDeleteType `json:"DeleteType"`
}
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
