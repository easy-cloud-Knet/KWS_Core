package domain

import (
	"errors"
	"fmt"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"github.com/google/uuid"
	"libvirt.org/go/libvirt"
)

type Domain interface {
	GetInfo() (*libvirt.DomainInfo, error)
	GetState() (libvirt.DomainState, int, error)
	GetUUID() ([]byte, error)
	GetGuestInfo(types libvirt.DomainGuestInfoTypes, flags uint32) (*libvirt.DomainGuestInfo, error)
}

type Connect interface {
	ListAllDomains(flags libvirt.ConnectListAllDomainsFlags) ([]libvirt.Domain, error)
}

type DataTypeHandler interface {
	GetInfo(Domain) error
}

type InstDataTypeHandler interface {
	GetAllinstInfo(Connect) error
}

type DataType uint

const (
	DomState DataType = iota
	BasicInfo
	GuestInfoUser
	GuestInfoOS
	GuestInfoFS
	GuestInfoDisk
)

type InstDataType uint

const (
	Vcpu_MaxMem InstDataType = iota
)

type State struct {
	DomainState libvirt.DomainState           `json:"currentState"`
	UUID        string                        `json:"UUID"`
	Users       []libvirt.DomainGuestInfoUser `json:"Guest Info,omitempty"`
}

type Info struct {
	State     libvirt.DomainState `json:"state"`
	MaxMem    uint64              `json:"maxmem"`
	Memory    uint64              `json:"memory"`
	NrVirtCpu uint                `json:"nrVirtCpu"`
	CpuTime   uint64              `json:"cpuTime"`
}

type AllInstInfo struct {
	Totalmaxmem uint64 `json:"totalmaxmem"`
	TotalVCpu   uint   `json:"totalVCpu"`
}

type Detail struct {
	DataHandle DataTypeHandler
	Domain     Domain
}

type InstDetail struct {
	AllInstDataHandle InstDataTypeHandler
}

func (DI *Info) GetInfo(domain Domain) error {
	info, err := domain.GetInfo()
	if err != nil {
		return virerr.ErrorGen(virerr.DomainStatusError, err)
	}
	DI.State = info.State
	DI.MaxMem = info.MaxMem
	DI.Memory = info.Memory
	DI.NrVirtCpu = info.NrVirtCpu
	DI.CpuTime = info.CpuTime
	return nil
}

func (DP *State) GetInfo(domain Domain) error {
	info, _, err := domain.GetState()
	if err != nil {
		return virerr.ErrorGen(virerr.DomainStatusError, err)
	}
	uuidBytes, err := domain.GetUUID()
	if err != nil {
		return virerr.ErrorGen(virerr.InvalidUUID, err)
	}
	uuidParsed, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return virerr.ErrorGen(virerr.InvalidUUID, err)
	}
	DP.DomainState = info
	DP.UUID = uuidParsed.String()
	userInfo, err := domain.GetGuestInfo(libvirt.DOMAIN_GUEST_INFO_USERS, 0)
	if err != nil {
		return virerr.ErrorGen(virerr.DomainStatusError, fmt.Errorf("error retreving guest info: %w", err))
	}
	DP.Users = userInfo.Users
	return nil
}

func (AII *AllInstInfo) GetAllinstInfo(conn Connect) error {
	domains, err := conn.ListAllDomains(0)
	if err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, fmt.Errorf("failed to list all domains: %w", err))
	}
	var totalMaxMem uint64
	var totalvCPU uint
	for _, dom := range domains {
		data, err := dom.GetInfo()
		if err != nil {
			dom.Free()
			continue
		}
		totalMaxMem += data.MaxMem
		totalvCPU += data.NrVirtCpu
		dom.Free()
	}
	AII.Totalmaxmem = totalMaxMem
	AII.TotalVCpu = totalvCPU
	return nil
}

func DetailFactory(handler DataTypeHandler, dom Domain) *Detail {
	return &Detail{
		DataHandle: handler,
		Domain:     dom,
	}
}

func DataTypeRouter(t DataType) (DataTypeHandler, error) {
	switch t {
	case DomState:
		return &State{}, nil
	case BasicInfo, GuestInfoUser, GuestInfoOS, GuestInfoFS, GuestInfoDisk:
		return &Info{}, nil
	}
	return nil, virerr.ErrorGen(virerr.InvalidParameter, errors.New("invalid flag for DataRoute entered"))
}

func InstDataTypeRouter(t InstDataType) (InstDataTypeHandler, error) {
	switch t {
	case Vcpu_MaxMem:
		return &AllInstInfo{}, nil
	}
	return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("unsupported type"))
}

func InstDetailFactory(handler InstDataTypeHandler, conn Connect) (*InstDetail, error) {
	if err := handler.GetAllinstInfo(conn); err != nil {
		return nil, err
	}
	return &InstDetail{AllInstDataHandle: handler}, nil
}
