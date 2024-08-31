package main

import (
	_ "log"

	"github.com/easy-cloud-Knet/KWS_Core.git/api"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/router"
)

func main() {
	libvirtInstance := router.LibvirtConnection()

	go api.Server(8080)
	router.MakeNewConnect(libvirtInstance)
	fmt.Println("working")	

	defer libvirtInstance.Close()



}
