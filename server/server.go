package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Core.git/api"
	syslogger "github.com/easy-cloud-Knet/KWS_Core.git/logger"
	"go.uber.org/zap"
)

func InitServer(portNum int, libvirtInst *api.InstHandler, logger zap.Logger) {
	logger.Sugar().Infof("Starting server on %d", portNum)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /createVM", libvirtInst.CreateVMFromBase)         //post
	mux.HandleFunc("GET /getStatusUUID", libvirtInst.ReturnStatusUUID)     //Get
	mux.HandleFunc("POST /forceShutDownUUID", libvirtInst.ForceShutDownVM) //POST
	mux.HandleFunc("POST /DeleteVM", libvirtInst.DeleteVM)                 //POST
	mux.HandleFunc("GET /getStatusHost", libvirtInst.ReturnStatusHost)     //Get
	mux.HandleFunc("GET /getInstAllInfo", libvirtInst.ReturnInstAllInfo)   //Get

	sysloggerHttp := syslogger.LoggerMiddleware(mux, &logger)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNum), sysloggerHttp))

}
