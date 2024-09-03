package conn

import (
	"fmt"
	"log"

	"libvirt.org/go/libvirt"
)




func (i *InstHandler)ReturnDomainNameList(flag libvirt.ConnectListAllDomainsFlags)([]*DomainInfo,error) {
	var Domains []*DomainInfo

	doms, err := i.LibvirtInst.ListAllDomains(flag)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {
		info, err := dom.GetInfo()
		
		if  err!=nil {
				log.Println(err)
		}
		
		DomInfo:= &DomainInfo{}
		DomInfo.State = info.State
		DomInfo.MaxMem = info.MaxMem
		DomInfo.Memory = info.Memory
		DomInfo.NrVirtCpu = info.NrVirtCpu
		DomInfo.CpuTime = info.CpuTime
		
		Domains:=append(Domains, DomInfo)
		Use(Domains)	
		dom.Free()
	}
	var erro1r error
	erro1r=nil
	return Domains,erro1r
}


func (i *InstHandler)LibvirtConnection(){
	i.InstMu.Lock()
	libvirtInst, err := libvirt.NewConnect("qemu:///system")
		if err != nil {
			panic(err)
}
	i.LibvirtInst = libvirtInst
	defer i.InstMu.Unlock()
}


func Use(vals ...interface{}){
	for _,val:= range vals{
	 _=val
	}
	

}


