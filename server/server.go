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


	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNum), nil))
}



