package status

import (
	"fmt"

	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domainList_status"
	"github.com/easy-cloud-Knet/KWS_Core/internal/domain"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"github.com/easy-cloud-Knet/KWS_Core/internal/host"
	"libvirt.org/go/libvirt"
)

type (
	DomainDataType      = domain.DataType
	DataTypeHandler     = domain.DataTypeHandler
	Domain              = domain.Domain
	Connect             = domain.Connect
	DomainDetail        = domain.Detail
	InstDataType        = domain.InstDataType
	InstDataTypeHandler = domain.InstDataTypeHandler
	InstDetail          = domain.InstDetail
	HostDataType        = host.DataType
	HostDataTypeHandler = host.DataTypeHandler
	HostDetail          = host.Detail
)

const (
	DomState      = domain.DomState
	BasicInfo     = domain.BasicInfo
	GuestInfoUser = domain.GuestInfoUser
	GuestInfoOS   = domain.GuestInfoOS
	GuestInfoFS   = domain.GuestInfoFS
	GuestInfoDisk = domain.GuestInfoDisk
	Vcpu_MaxMem   = domain.Vcpu_MaxMem
)

func DataTypeRouter(t DomainDataType) (DataTypeHandler, error) {
	return domain.DataTypeRouter(t)
}

func DomainDetailFactory(handler DataTypeHandler, dom Domain) *DomainDetail {
	return domain.DetailFactory(handler, dom)
}

func HostDataTypeRouter(t HostDataType) (HostDataTypeHandler, error) {
	return host.DataTypeRouter(t)
}

func HostInfoHandler(handler HostDataTypeHandler, s *domStatus.DomainListStatus) (*HostDetail, error) {
	return host.InfoHandler(handler, s)
}

func InstDataTypeRouter(t InstDataType) (InstDataTypeHandler, error) {
	return domain.InstDataTypeRouter(t)
}

func InstDetailFactory(handler InstDataTypeHandler, conn Connect) (*InstDetail, error) {
	return domain.InstDetailFactory(handler, conn)
}

type DomainStateInfo struct {
	UUID        string              `json:"UUID"`
	DomainState libvirt.DomainState `json:"currentState"`
}

func ListAllDomainStates(conn Connect) ([]DomainStateInfo, error) {
	domains, err := conn.ListAllDomains(0)
	if err != nil {
		return nil, virerr.ErrorGen(virerr.DomainStatusError, fmt.Errorf("failed to list domains: %w", err))
	}
	defer func() {
		for i := range domains {
			domains[i].Free()
		}
	}()
	result := make([]DomainStateInfo, 0, len(domains))
	for i := range domains {
		uuid, err := domains[i].GetUUIDString()
		if err != nil {
			continue
		}
		state, _, err := domains[i].GetState()
		if err != nil {
			continue
		}
		result = append(result, DomainStateInfo{UUID: uuid, DomainState: state})
	}
	return result, nil
}
