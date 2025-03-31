package api

import (
	"os"

	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)


 

func (i *InstHandler)LibvirtConnection(){
	libvirtInst, err := libvirt.NewConnect("qemu:///system")
		if err != nil {
			i.Logger.Panic("innitial connection for libvirt daemon failed", zap.Int("pid", os.Getegid()))
}
	i.LibvirtInst = libvirtInst
	i.Logger.Info("Libvirt Coonnection succefully done.", zap.Int("pid",os.Getegid()))
	
	defer i.Logger.Sync()
}


 

