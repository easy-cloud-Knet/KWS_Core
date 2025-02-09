package network


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
	Ips []string `json:"ips"`
	NetType NetType  `json:"NetType"`
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