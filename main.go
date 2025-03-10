package main

import (
	_ "context"
	_ "log"
	"runtime/debug"

	domCon "github.com/easy-cloud-Knet/KWS_Core.git/api/conn/DomCon"
	syslogger "github.com/easy-cloud-Knet/KWS_Core.git/api/logger"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/service"
	"github.com/easy-cloud-Knet/KWS_Core.git/server"
)

func main() {
	debug.SetTraceback("none")
	logger := syslogger.InitialLogger()

	domListCon := domCon.DomListConGen()

	libvirtInst := service.InstHandler{
		DomainControl: domListCon,
		Logger:        logger,
	}

	libvirtInst.LibvirtConnection()
	libvirtInst.DomainControl.RetrieveAllDomain(libvirtInst.LibvirtInst, logger)

	go server.InitServer(8080, &libvirtInst, *logger)

	defer func() {

		logger.Info("Shutting down gracefully...") // 종료 시 로깅
		logger.Sync()
		libvirtInst.LibvirtInst.Close()
	}()
	select {}
}
