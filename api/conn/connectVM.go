package conn

import (
	"fmt"

	libvirt "libvirt.org/go/libvirt"
)

type currentInst struct{
	LibvirtInst *libvirt.Connect
}

var currLibvirtInst currentInst

func ActiveDomain(libvirtInst *libvirt.Connect) {
	
	doms, err := libvirtInst.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
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


func LibvirtConnection() *libvirt.Connect{
	libvirtInst, err := libvirt.NewConnect("qemu:///system")
		if err != nil {
			panic(err)
}
	return libvirtInst
}

func SetLibvirtInst(libvirtInst *libvirt.Connect){
	currLibvirtInst.LibvirtInst = libvirtInst
}
func ReturnLibvirtInst() currentInst{
	return currLibvirtInst
}
