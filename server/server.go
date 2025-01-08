package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)


func InitServer(portNum int, libvirtInst *conn.InstHandler, domain *conn.Domain){
	
	http.HandleFunc("GET /getStatus", libvirtInst.ReturnStatus)
	http.HandleFunc("POST /createVM", libvirtInst.CreateVM)
	http.HandleFunc("Get /getStatusUUID", libvirtInst.ReturnStatusUUID)
	http.HandleFunc("Get /ForceShutDownUUID", libvirtInst.ForceShutDownVM)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNum), nil))
}



