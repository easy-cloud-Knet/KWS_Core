package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/service"
	timeCal "github.com/easy-cloud-Knet/KWS_Core.git/server/test.go"
)

func InitServer(portNum int, libvirtInst *service.InstHandler) {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /getStatus", libvirtInst.ReturnDomainByStatus)    //get
	mux.HandleFunc("POST /createVM", libvirtInst.CreateVMLocal)                 //post
	mux.HandleFunc("GET /getStatusUUID", libvirtInst.ReturnStatusUUID)    //Get
	mux.HandleFunc("POST /forceShutDownUUID", libvirtInst.ForceShutDownVM) //Get
	mux.HandleFunc("POST /DeleteVM", libvirtInst.DeleteVM)                 //Get
	mux.HandleFunc("GET /getStatusHost", libvirtInst.ReturnStatusHost)    //Get

	timeCalulatorHTTP := timeCal.TimeLogging(mux) 


	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNum), timeCalulatorHTTP))
}
