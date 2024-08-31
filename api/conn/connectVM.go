package conn

import (
	"fmt"
	libvirt "libvirt.org/go/libvirt"
)
func ActiveDomain(libvirtInst *libvirt.Connect) {
	
	doms, err := libvirtInst.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {
		name, err := dom.GetName()
		if err == nil {
			fmt.Printf("  %s\n", name)
		}
		dom.Free()
	}
	defer libvirtInst.Close()
}


func LibvirtConnection() *libvirt.Connect{
	
	libvirtInst, err := libvirt.NewConnect("qemu:///system")
		if err != nil {
			panic(err)
		
}
	


	return libvirtInst
}
