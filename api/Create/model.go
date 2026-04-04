package create

import (
	network "github.com/easy-cloud-Knet/KWS_Core/internal/net"
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
)

type CreateVMRequest struct {
	DomName      string                 `json:"domName"`
	UUID         string                 `json:"uuid"`
	OS           string                 `json:"os"`
	HardwardInfo vmtypes.HardwareInfo   `json:"HWInfo"`
	NetConf      network.NetDefine      `json:"network"`
	Users        []vmtypes.User_info_VM `json:"users"`
	SDNUUID      string                 `json:"sdnUUID"`
	MacAddr      string                 `json:"macAddr"`
}

func (r *CreateVMRequest) toVMInitInfo() *vmtypes.VM_Init_Info {
	return &vmtypes.VM_Init_Info{
		DomName:      r.DomName,
		UUID:         r.UUID,
		OS:           r.OS,
		HardwardInfo: r.HardwardInfo,
		NetConf:      r.NetConf,
		Users:        r.Users,
		SDNUUID:      r.SDNUUID,
		MacAddr:      r.MacAddr,
	}
}

type DomainBootRequest struct {
	UUID string `json:"UUID"`
}
