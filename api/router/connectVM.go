package router

import (
	"fmt"

	libvirt "libvirt.org/libvirt-go"
)

func MakeNewConnect(libvirt.Conn) {
	
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
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
	defer conn.Close()
}

func LibvirtConnection() *libvirt.Connect{
	
	conn, err := libvirt.NewConnect("qemu:///system")
		if err != nil {
			panic(err)
		
}
	defer conn

}
