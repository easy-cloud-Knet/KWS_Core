package main

import (
	_ "context"
	"fmt"
	_ "log"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/service"
	"github.com/easy-cloud-Knet/KWS_Core.git/server"
)

func main() {
	var libvirtInst service.InstHandler
	libvirtInst.LibvirtConnection()

	go server.InitServer(8080, &libvirtInst)
	fmt.Println("working")	

	defer  libvirtInst.LibvirtInst.Close()

	select {}
}
