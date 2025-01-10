package main

import (
	_ "context"
	"fmt"
	_ "log"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"github.com/easy-cloud-Knet/KWS_Core.git/server"
)

func main() {
	var libvirtInst conn.InstHandler
	var domain conn.Domain
	libvirtInst.LibvirtConnection()

	go server.InitServer(8080, &libvirtInst, &domain)
	fmt.Println("working")	

	defer  libvirtInst.LibvirtInst.Close()

	select {}
}
