package server

import (
	"net/http"
	"strconv"

	control "github.com/easy-cloud-Knet/KWS_Core/api/Control"
	create "github.com/easy-cloud-Knet/KWS_Core/api/Create"
	snapshot "github.com/easy-cloud-Knet/KWS_Core/api/Snapshot"
	apistatus "github.com/easy-cloud-Knet/KWS_Core/api/Status"
	"github.com/easy-cloud-Knet/KWS_Core/api/metric"
	"github.com/easy-cloud-Knet/KWS_Core/internal/server/middleware"
	"go.uber.org/zap"
)

type Handlers struct {
	IsConnected func() bool
	Control     *control.Handler
	Metric      *metric.Handler
	Create      *create.Handler
	Snapshot    *snapshot.Handler
	Status      *apistatus.Handler
}

func InitServer(portNum int, h Handlers, logger *zap.Logger) {
	logger.Sugar().Infof("Starting server on %d", portNum)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /BOOTVM", h.Create.BootVM)
	mux.HandleFunc("POST /createVM", h.Create.CreateVMFromBase)
	mux.HandleFunc("GET /getStatusUUID", h.Status.ReturnStatusUUID)
	mux.HandleFunc("POST /forceShutDownUUID", h.Control.ForceShutDownVM)
	mux.HandleFunc("POST /DeleteVM", h.Control.DeleteVM)
	mux.HandleFunc("GET /getStatusHost", h.Status.ReturnStatusHost)
	mux.HandleFunc("GET /getInstAllInfo", h.Status.ReturnInstAllInfo)
	mux.HandleFunc("GET /getAllUUIDs", h.Status.ReturnAllUUIDs)
	mux.HandleFunc("GET /getAll-uuidstatusList", h.Status.ReturnAllDomainStates)

	mux.HandleFunc("POST /CreateSnapshot", h.Snapshot.CreateSnapshot)
	mux.HandleFunc("POST /CreateExternalSnapshot", h.Snapshot.CreateExternalSnapshot)
	mux.HandleFunc("GET /ListSnapshots", h.Snapshot.ListSnapshots)
	mux.HandleFunc("GET /ListExternalSnapshots", h.Snapshot.ListExternalSnapshots)
	mux.HandleFunc("POST /RevertSnapshot", h.Snapshot.RevertSnapshot)
	mux.HandleFunc("POST /RevertExternalSnapshot", h.Snapshot.RevertExternalSnapshot)
	mux.HandleFunc("POST /MergeExternalSnapshot", h.Snapshot.MergeExternalSnapshot)
	mux.HandleFunc("POST /DeleteSnapshot", h.Snapshot.DeleteSnapshot)

	mux.HandleFunc("GET /metrics", h.Metric.DefaultMetric().ServeHTTP)

	libvirtHandler := middleware.LibvirtMiddleware(h.IsConnected, logger)(mux)
	sysloggerHttp := middleware.LoggerMiddleware(libvirtHandler, logger)

	if err := http.ListenAndServe(":"+strconv.Itoa(portNum), sysloggerHttp); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
