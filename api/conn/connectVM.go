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
		uuid, err := dom.GetUUIDString()
		if err!= nil{
			log.Panicln(err)
		}
		DomInfo:= &DomainInfo{
				State :info.State,
				MaxMem :info.MaxMem,
				Memory : info.Memory,
				NrVirtCpu :info.NrVirtCpu,
				CpuTime :info.CpuTime,
				UUID : uuid,
		}
		
		Domains=append(Domains,DomInfo)
		dom.Free()
	}
	var erro1r error
	erro1r=nil
	Use(Domains)
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


