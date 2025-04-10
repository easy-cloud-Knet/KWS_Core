package status

import (
	"errors"
	"log"

	domCon "github.com/easy-cloud-Knet/KWS_Core.git/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core.git/error"
	"github.com/google/uuid"
	"libvirt.org/go/libvirt"
)

func (DI *DomainInfo) GetInfo(domain *domCon.Domain) error {
	info, err := domain.Domain.GetInfo()
	if err != nil {
		return virerr.ErrorGen(virerr.DomainStatusError, err)
	}
	DI.State = info.State
	DI.MaxMem = info.MaxMem
	DI.Memory = info.Memory
	DI.NrVirtCpu = info.NrVirtCpu
	DI.CpuTime = info.CpuTime
	//basic info can be added
	return nil
}

func (DP *DomainState) GetInfo(domain *domCon.Domain) error {
	info, _, err := domain.Domain.GetState()
	//searching for coresponding second parameter, "Reason"
	if err != nil {
		return virerr.ErrorGen(virerr.DomainStatusError, err)
	}

	uuidBytes, err := domain.Domain.GetUUID()
	if err != nil {
		return virerr.ErrorGen(virerr.InvalidUUID, err)
	}
	uuidParsed, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return virerr.ErrorGen(virerr.InvalidUUID, err)
	}

	DP.DomainState = info
	DP.UUID = string(uuidParsed.String())
	userInfo, err := domain.Domain.GetGuestInfo(libvirt.DOMAIN_GUEST_INFO_USERS, 0)
	if err != nil {
		log.Println("error retreving guest info")
		return err
	}
	DP.Users = userInfo.Users
	return nil
}

func DomainDetailFactory(Handler DataTypeHandler, dom *domCon.Domain) *DomainDetail {
	return &DomainDetail{
		DataHandle: Handler,
		Domain:     dom,
	}
}

func DataTypeRouter(types DomainDataType) (DataTypeHandler, error) {
	switch types {
	case DomState:
		return &DomainState{}, nil
	case BasicInfo:
		return &DomainInfo{}, nil
	case GuestInfoUser:
		return &DomainInfo{}, nil
	case GuestInfoOS:
		return &DomainInfo{}, nil
	case GuestInfoFS:
		return &DomainInfo{}, nil
	case GuestInfoDisk:
		return &DomainInfo{}, nil

	}
	return nil, virerr.ErrorGen(virerr.InvalidParameter, errors.New("invalid flag for DataRoute entereed "))
}
