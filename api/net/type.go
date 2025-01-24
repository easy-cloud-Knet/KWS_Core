package net


type NetType uint

// type currentNetdevicesHost struct{
// 	AllocatedNetInt []NetInterface 
// 	UnAllocatedNetInt []NetInterface
// }

const (
	Bridge = iota
	Nat

)
type NetDefine struct{
	NetType NetType  `json:"NetType"`
	Ips []string `json:"ips"`
}
// gonna replace fields in VM_Init_Info
//structure,need to modify parsor when implement this 

type NetInterface struct{
	PortNumber uint8 
	NetType NetType `json:"NetType"`
}

type VnetInterface struct{
}

//currently building for advanced 