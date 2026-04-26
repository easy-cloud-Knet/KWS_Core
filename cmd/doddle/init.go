package main

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	control "github.com/easy-cloud-Knet/KWS_Core/api/Control"
	create "github.com/easy-cloud-Knet/KWS_Core/api/Create"
	apistatus "github.com/easy-cloud-Knet/KWS_Core/api/Status"
	snapshot "github.com/easy-cloud-Knet/KWS_Core/api/Snapshot"
	libvirtconn "github.com/easy-cloud-Knet/KWS_Core/internal/libvirt"
	syslogger "github.com/easy-cloud-Knet/KWS_Core/internal/logger"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type App struct {
	LibvirtConn     *libvirt.Connect
	ControlHandler  *control.Handler
	CreateHandler   *create.Handler
	SnapshotHandler *snapshot.Handler
	StatusHandler   *apistatus.Handler
	Logger          *zap.Logger
}

func (a *App) IsConnected() bool {
	return libvirtconn.IsAlive(a.LibvirtConn)
}

func initApp() *App {
	logger := syslogger.InitialLogger()

	conn, err := libvirtconn.Connect(logger)
	if err != nil {
		logger.Fatal("initial connection for libvirt daemon failed", zap.Error(err))
	}

	domListCon := domCon.DomListConGen()
	domListCon.SetLibvirt(conn)
	domListCon.DomainListStatus.Update()

	if err := domListCon.RetrieveAllDomain(logger); err != nil {
		logger.Fatal("failed to retrieve domains on startup", zap.Error(err))
	}

	return &App{
		LibvirtConn:     conn,
		ControlHandler:  control.NewHandler(domListCon, logger),
		CreateHandler:   create.NewHandler(domListCon, conn, logger),
		SnapshotHandler: snapshot.NewHandler(domListCon, logger),
		StatusHandler:   apistatus.NewHandler(conn, domListCon, logger),
		Logger:          logger,
	}
}

func (a *App) Shutdown() {
	a.Logger.Info("Shutting down gracefully...")
	a.Logger.Sync()
	a.LibvirtConn.Close()
}
