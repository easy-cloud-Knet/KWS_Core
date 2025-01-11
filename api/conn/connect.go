package conn

import (
	"libvirt.org/go/libvirt"
)


func (i *InstHandler)LibvirtConnection(){
	libvirtInst, err := libvirt.NewConnect("qemu:///system")
		if err != nil {
			panic(err)
}
	i.LibvirtInst = libvirtInst
}


 
