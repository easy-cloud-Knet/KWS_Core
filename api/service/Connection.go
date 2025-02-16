package service

import (
	"libvirt.org/go/libvirt"
)


 

func (i *InstHandler)LibvirtConnection(){
	libvirtInst, err := libvirt.NewConnect("qemu:///system")
		if err != nil {
			panic(err)
}
	i.LibvirtInst = libvirtInst
	i.Logger.Info("Libvirt Coonnection succefully done.")
	
	defer i.Logger.Sync()
}


 

