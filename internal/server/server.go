package server

import (
	"net/http"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Core/api"
	control "github.com/easy-cloud-Knet/KWS_Core/api/Control"
	create "github.com/easy-cloud-Knet/KWS_Core/api/Create"
	snapshot "github.com/easy-cloud-Knet/KWS_Core/api/Snapshot"
	"github.com/easy-cloud-Knet/KWS_Core/internal/server/middleware"
	"go.uber.org/zap"
)

func InitServer(portNum int, libvirtInst *api.InstHandler, controlHandler *control.Handler, createHandler *create.Handler, snapshotHandler *snapshot.Handler, logger *zap.Logger) {
	logger.Sugar().Infof("Starting server on %d", portNum)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /BOOTVM", createHandler.BootVM)                      //post
	mux.HandleFunc("POST /createVM", createHandler.CreateVMFromBase)          //post
	mux.HandleFunc("GET /getStatusUUID", libvirtInst.ReturnStatusUUID)        //Get
	mux.HandleFunc("POST /forceShutDownUUID", controlHandler.ForceShutDownVM) //POST
	mux.HandleFunc("POST /DeleteVM", controlHandler.DeleteVM)                 //POST
	mux.HandleFunc("GET /getStatusHost", libvirtInst.ReturnStatusHost)        //Get
	mux.HandleFunc("GET /getInstAllInfo", libvirtInst.ReturnInstAllInfo)      //Get
	mux.HandleFunc("GET /getAllUUIDs", libvirtInst.ReturnAllUUIDs)            //Get
	mux.HandleFunc("GET /getAll-uuidstatusList", libvirtInst.ReturnAllDomainStates)

	// Snapshot operations
	mux.HandleFunc("POST /CreateSnapshot", snapshotHandler.CreateSnapshot)
	mux.HandleFunc("POST /CreateExternalSnapshot", snapshotHandler.CreateExternalSnapshot)
	mux.HandleFunc("GET /ListSnapshots", snapshotHandler.ListSnapshots)
	mux.HandleFunc("GET /ListExternalSnapshots", snapshotHandler.ListExternalSnapshots)
	mux.HandleFunc("POST /RevertSnapshot", snapshotHandler.RevertSnapshot)
	mux.HandleFunc("POST /RevertExternalSnapshot", snapshotHandler.RevertExternalSnapshot)
	mux.HandleFunc("POST /MergeExternalSnapshot", snapshotHandler.MergeExternalSnapshot)
	mux.HandleFunc("POST /DeleteSnapshot", snapshotHandler.DeleteSnapshot)

	libvirtHandler := middleware.LibvirtMiddleware(libvirtInst.IsConnected, logger)(mux)
	sysloggerHttp := middleware.LoggerMiddleware(libvirtHandler, logger)

	if err := http.ListenAndServe(":"+strconv.Itoa(portNum), sysloggerHttp); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}

}
