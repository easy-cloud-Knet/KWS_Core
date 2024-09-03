package conn

import (
	"fmt"

	"libvirt.org/go/libvirt"
)




func (i *InstHandler)ReturnDomainNameList(flag libvirt.ConnectListAllDomainsFlags) {

	doms, err := i.LibvirtInst.ListAllDomains(flag)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {
		name, err := dom.GetName()
		if err == nil {
			fmt.Printf("%s\n", name)
		}
		dom.Free()
	}
	
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


