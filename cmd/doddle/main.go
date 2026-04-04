package main

import (
	_ "context"
	_ "log"
	"runtime/debug"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core/api"
	control "github.com/easy-cloud-Knet/KWS_Core/api/Control"
	"github.com/easy-cloud-Knet/KWS_Core/internal/config"
	syslogger "github.com/easy-cloud-Knet/KWS_Core/internal/logger"
	"github.com/easy-cloud-Knet/KWS_Core/internal/server"
	"go.uber.org/zap"
)

func main() {
	debug.SetTraceback("none")
	logger := syslogger.InitialLogger()

	domListCon := domCon.DomListConGen()

	libvirtInst := api.InstHandler{
		DomainControl: domListCon,
		Logger:        logger,
	}

	libvirtInst.LibvirtConnection()
	libvirtInst.DomainControl.SetLibvirt(libvirtInst.LibvirtInst)
	libvirtInst.DomainControl.DomainListStatus.Update()
	if err := libvirtInst.DomainControl.RetrieveAllDomain(logger); err != nil {
		logger.Fatal("failed to retrieve domains on startup", zap.Error(err))
	}

	controlHandler := control.NewHandler(domListCon, logger)

	go server.InitServer(config.ServerPort, &libvirtInst, controlHandler, logger)
	defer func() {
		logger.Info("Shutting down gracefully...") // 종료 시 로깅
		logger.Sync()
		libvirtInst.LibvirtInst.Close()
	}()
	select {}
}
