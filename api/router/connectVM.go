package router

import (
	"fmt"
	libvirt "libvirt.org/go/libvirt"
)
func Makenewconnect(conn *libvirt.Connect) {
	
	doms, err := conn.listalldomains(libvirt.connect_list_domains_active)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {
		name, err := dom.getname()
		if err == nil {
			fmt.Printf("  %s\n", name)
		}
		dom.free()
	}
	defer conn.close()
}


func LibvirtConnection() *libvirt.Connect{
	
	conn, err := libvirt.NewConnect("qemu:///system")
		if err != nil {
			panic(err)
		
}
	


	return conn
}
