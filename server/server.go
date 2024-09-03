package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)


func InitServer(portNum int, libvirtInst *conn.InstHandler){
	
	http.HandleFunc("GET /getStatus", libvirtInst.ReturnStatus)
	// http.HandleFunc("POST /createVM")


	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNum), nil))
}



