package main

import (
	_ "log"
	"fmt"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/server"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)

func main() {
	a:=make (chan int)
	
	libvirtInst := conn.LibvirtConnection()

	go server.InitServer(8080, libvirtInst)
	fmt.Println("working")	

	defer libvirtInst.Close()


	<-a
}
