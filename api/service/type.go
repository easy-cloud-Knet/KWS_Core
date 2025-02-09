package service

import (
	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"libvirt.org/go/libvirt"
)

type InstHandler struct {
	LibvirtInst conn.LibvirtInst
}

type InstHandle interface {
	LibvirtConnection()
}

//// income api Structures

type DeleteDomain struct {
	UUID         string                `json:"UUID"`
	DeletionType conn.DomainDeleteType `json:"DeleteType"`
}
type ShutDownDomain struct {
	UUID string `json:"UUID"`
}
type StartDomain struct {
	UUID string `json:"UUID"`
}

type ReturnDomainFromStatus struct {
	DataType conn.DomainDataType                `json:"dataType"`
	Status   libvirt.ConnectListAllDomainsFlags `json:"Flag"`
}

type ReturnDomainFromUUID struct {
	DataType conn.DomainDataType `json:"dataType"`
	UUID     string              `json:"UUID"`
}

// host
type ReturnHostFromStatus struct {
	HostDataType conn.HostDataType `json:"host_dataType"`
}

//

// type ConnectListAllDomainsFlags uint
//     const (
//     CONNECT_LIST_DOMAINS_ACTIVE         = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_ACTIVE)
//     CONNECT_LIST_DOMAINS_INACTIVE       = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_INACTIVE)
//     CONNECT_LIST_DOMAINS_PERSISTENT     = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_PERSISTENT)
//     CONNECT_LIST_DOMAINS_TRANSIENT      = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_TRANSIENT)
//     CONNECT_LIST_DOMAINS_RUNNING        = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_RUNNING)
//     CONNECT_LIST_DOMAINS_PAUSED         = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_PAUSED)
//     CONNECT_LIST_DOMAINS_SHUTOFF        = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_SHUTOFF)
//     CONNECT_LIST_DOMAINS_OTHER          = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_OTHER)
//     CONNECT_LIST_DOMAINS_MANAGEDSAVE    = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_MANAGEDSAVE)
//     CONNECT_LIST_DOMAINS_NO_MANAGEDSAVE = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_NO_MANAGEDSAVE)
//     CONNECT_LIST_DOMAINS_AUTOSTART      = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_AUTOSTART)
//     CONNECT_LIST_DOMAINS_NO_AUTOSTART   = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_NO_AUTOSTART)
//     CONNECT_LIST_DOMAINS_HAS_SNAPSHOT   = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_HAS_SNAPSHOT)
//     CONNECT_LIST_DOMAINS_NO_SNAPSHOT    = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_NO_SNAPSHOT)
//     CONNECT_LIST_DOMAINS_HAS_CHECKPOINT = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_HAS_CHECKPOINT)
//     CONNECT_LIST_DOMAINS_NO_CHECKPOINT  = ConnectListAllDomainsFlags(C.VIR_CONNECT_LIST_DOMAINS_NO_CHECKPOINT)
// )

// const (
// 	DOMAIN_NOSTATE     = DomainState(C.VIR_DOMAIN_NOSTATE)
// 	DOMAIN_RUNNING     = DomainState(C.VIR_DOMAIN_RUNNING)
// 	DOMAIN_BLOCKED     = DomainState(C.VIR_DOMAIN_BLOCKED)
// 	DOMAIN_PAUSED      = DomainState(C.VIR_DOMAIN_PAUSED)
// 	DOMAIN_SHUTDOWN    = DomainState(C.VIR_DOMAIN_SHUTDOWN)
// 	DOMAIN_CRASHED     = DomainState(C.VIR_DOMAIN_CRASHED)
// 	DOMAIN_PMSUSPENDED = DomainState(C.VIR_DOMAIN_PMSUSPENDED)
// 	DOMAIN_SHUTOFF     = DomainState(C.VIR_DOMAIN_SHUTOFF)
// )
