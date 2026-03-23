package vmtypes

import network "github.com/easy-cloud-Knet/KWS_Core/net"

type VM_Init_Info struct {
	DomName      string            `json:"domName"`
	UUID         string            `json:"uuid"`
	OS           string            `json:"os"`
	HardwardInfo HardwareInfo      `json:"HWInfo"`
	NetConf      network.NetDefine `json:"network"`
	Users        []User_info_VM    `json:"users"`
	SDNUUID      string            `json:"sdnUUID"`
	MacAddr      string            `json:"macAddr"`
}

type HardwareInfo struct {
	CPU    int `json:"cpu"`
	Memory int `json:"memory"`
	Disk   int `json:"disk"`
}

type User_info_VM struct {
	Name                string   `json:"name"`
	Groups              string   `json:"groups"`
	PassWord            string   `json:"passWord"`
	Ssh_authorized_keys []string `json:"ssh"`
}
