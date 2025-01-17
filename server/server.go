package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/service"
)


func InitServer(portNum int, libvirtInst *service.InstHandler){
	
	http.HandleFunc("/getStatus", libvirtInst.ReturnDomainByStatus) //get
	http.HandleFunc("/createVM", libvirtInst.CreateVM) //post
	http.HandleFunc("/getStatusUUID", libvirtInst.ReturnStatusUUID) //Get
	http.HandleFunc("/forceShutDownUUID", libvirtInst.ForceShutDownVM) //Get
	http.HandleFunc("/DeleteVM", libvirtInst.DeleteVM) //Get


	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNum), nil))
}



