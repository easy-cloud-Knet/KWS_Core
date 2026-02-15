package server

import (
	"net/http"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Core/api"
	syslogger "github.com/easy-cloud-Knet/KWS_Core/logger"
	"go.uber.org/zap"
)

func InitServer(portNum int, libvirtInst *api.InstHandler, logger *zap.Logger) {
	logger.Sugar().Infof("Starting server on %d", portNum)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /BOOTVM", libvirtInst.BootVM)                     //post
	mux.HandleFunc("POST /createVM", libvirtInst.CreateVMFromBase)         //post
	mux.HandleFunc("GET /getStatusUUID", libvirtInst.ReturnStatusUUID)     //Get
	mux.HandleFunc("POST /forceShutDownUUID", libvirtInst.ForceShutDownVM) //POST
	mux.HandleFunc("POST /DeleteVM", libvirtInst.DeleteVM)                 //POST
	mux.HandleFunc("GET /getStatusHost", libvirtInst.ReturnStatusHost)     //Get
	mux.HandleFunc("GET /getInstAllInfo", libvirtInst.ReturnInstAllInfo)   //Get
	mux.HandleFunc("GET /getAllUUIDs", libvirtInst.ReturnAllUUIDs)         //Get
	mux.HandleFunc("GET /getAll-uuidstatusList", libvirtInst.ReturnAllDomainStates)

	// Snapshot operations
	mux.HandleFunc("POST /CreateSnapshot", libvirtInst.CreateSnapshot)
	mux.HandleFunc("GET /ListSnapshots", libvirtInst.ListSnapshots)
	mux.HandleFunc("POST /RevertSnapshot", libvirtInst.RevertSnapshot)
	mux.HandleFunc("POST /DeleteSnapshot", libvirtInst.DeleteSnapshot)

	sysloggerHttp := syslogger.LoggerMiddleware(mux, logger)

	if err := http.ListenAndServe(":"+strconv.Itoa(portNum), sysloggerHttp); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}

}
