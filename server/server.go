package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)


func InitServer(portNum int, libvirtInst *conn.InstHandler, domain *conn.Domain){
	
	http.HandleFunc("/getStatus", libvirtInst.ReturnDomainByStatus) //get
	http.HandleFunc("/createVM", libvirtInst.CreateVM) //post
	http.HandleFunc("/getStatusUUID", libvirtInst.ReturnStatusUUID) //Get
	http.HandleFunc("/forceShutDownUUID", libvirtInst.ForceShutDownVM) //Get

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNum), nil))
}



