package main

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core/api"
	control "github.com/easy-cloud-Knet/KWS_Core/api/Control"
	create "github.com/easy-cloud-Knet/KWS_Core/api/Create"
	syslogger "github.com/easy-cloud-Knet/KWS_Core/internal/logger"
	"go.uber.org/zap"
)

type App struct {
	InstHandler   *api.InstHandler
	ControlHandler *control.Handler
	CreateHandler  *create.Handler
	Logger         *zap.Logger
}

func initApp() *App {
	logger := syslogger.InitialLogger()

	domListCon := domCon.DomListConGen()

	inst := &api.InstHandler{
		DomainControl: domListCon,
		Logger:        logger,
	}

	inst.LibvirtConnection()
	inst.DomainControl.SetLibvirt(inst.LibvirtInst)
	inst.DomainControl.DomainListStatus.Update()

	if err := inst.DomainControl.RetrieveAllDomain(logger); err != nil {
		logger.Fatal("failed to retrieve domains on startup", zap.Error(err))
	}

	return &App{
		InstHandler:    inst,
		ControlHandler: control.NewHandler(domListCon, logger),
		CreateHandler:  create.NewHandler(domListCon, inst.LibvirtInst, logger),
		Logger:         logger,
	}
}

func (a *App) Shutdown() {
	a.Logger.Info("Shutting down gracefully...")
	a.Logger.Sync()
	a.InstHandler.LibvirtInst.Close()
}
