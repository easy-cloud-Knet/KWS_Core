package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"libvirt.org/go/libvirt"
)







func InitServer(portNum int, libvirtInst *libvirt.Connect){
	
	http.HandleFunc("/getStatus",status)



	http.ListenAndServe(":"+strconv.Itoa(portNum), nil)
}


func status(w http.ResponseWriter,r * http.Request){
	fmt.Println("getStatus request income")
	conn.ActiveDomain(conn.ReturnLibvirtInst().LibvirtInst)
}

