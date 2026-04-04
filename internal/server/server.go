package server

import (
	"net/http"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Core/api"
	create "github.com/easy-cloud-Knet/KWS_Core/api/Create"
	"github.com/easy-cloud-Knet/KWS_Core/internal/server/middleware"
	"go.uber.org/zap"
)

func InitServer(portNum int, libvirtInst *api.InstHandler, createHandler *create.Handler, logger *zap.Logger) {
	logger.Sugar().Infof("Starting server on %d", portNum)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /BOOTVM", createHandler.BootVM)                   //post
	mux.HandleFunc("POST /createVM", createHandler.CreateVMFromBase)       //post
	mux.HandleFunc("GET /getStatusUUID", libvirtInst.ReturnStatusUUID)     //Get
	mux.HandleFunc("POST /forceShutDownUUID", libvirtInst.ForceShutDownVM) //POST
	mux.HandleFunc("POST /DeleteVM", libvirtInst.DeleteVM)                 //POST
	mux.HandleFunc("GET /getStatusHost", libvirtInst.ReturnStatusHost)     //Get
	mux.HandleFunc("GET /getInstAllInfo", libvirtInst.ReturnInstAllInfo)   //Get
	mux.HandleFunc("GET /getAllUUIDs", libvirtInst.ReturnAllUUIDs)         //Get
	mux.HandleFunc("GET /getAll-uuidstatusList", libvirtInst.ReturnAllDomainStates)

	// Snapshot operations
	mux.HandleFunc("POST /CreateSnapshot", libvirtInst.CreateSnapshot)
	mux.HandleFunc("POST /CreateExternalSnapshot", libvirtInst.CreateExternalSnapshot)
	mux.HandleFunc("GET /ListSnapshots", libvirtInst.ListSnapshots)
	mux.HandleFunc("GET /ListExternalSnapshots", libvirtInst.ListExternalSnapshots)
	mux.HandleFunc("POST /RevertSnapshot", libvirtInst.RevertSnapshot)
	mux.HandleFunc("POST /RevertExternalSnapshot", libvirtInst.RevertExternalSnapshot)
	mux.HandleFunc("POST /MergeExternalSnapshot", libvirtInst.MergeExternalSnapshot)
	mux.HandleFunc("POST /DeleteSnapshot", libvirtInst.DeleteSnapshot)

	libvirtHandler := middleware.LibvirtMiddleware(libvirtInst.IsConnected, logger)(mux)
	sysloggerHttp := middleware.LoggerMiddleware(libvirtHandler, logger)

	if err := http.ListenAndServe(":"+strconv.Itoa(portNum), sysloggerHttp); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}

}
