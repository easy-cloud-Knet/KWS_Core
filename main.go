package main

import (
	"fmt"
	_ "log"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/server"
)

func main() {
	libvirtInst :=conn.LibvirtConnection()
	conn.SetLibvirtInst(libvirtInst)

	go server.InitServer(8080, libvirtInst)
	fmt.Println("working")	

	defer  libvirtInst.Close()


	select {}
}
