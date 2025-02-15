package main

import (
	_ "context"
	"fmt"
	_ "log"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/service"
	"github.com/easy-cloud-Knet/KWS_Core.git/server"
)

func main() {
	domListCon := conn.DomListConGen()

	libvirtInst := service.InstHandler{
		DomainControl:domListCon,
	}
	libvirtInst.LibvirtConnection()
	libvirtInst.DomainControl.RetreiveAllDomain(libvirtInst.LibvirtInst)

	go server.InitServer(8080, &libvirtInst)
	fmt.Println("working")

	defer libvirtInst.LibvirtInst.Close()

	select {}
}
