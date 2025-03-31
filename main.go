package main

import (
	_ "context"
	_ "log"
	"runtime/debug"

	domCon "github.com/easy-cloud-Knet/KWS_Core.git/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core.git/api"
	syslogger "github.com/easy-cloud-Knet/KWS_Core.git/logger"
	"github.com/easy-cloud-Knet/KWS_Core.git/server"
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
	libvirtInst.DomainControl.RetrieveAllDomain(libvirtInst.LibvirtInst, logger)

	go server.InitServer(8080, &libvirtInst, *logger)

	defer func() {

		logger.Info("Shutting down gracefully...") // 종료 시 로깅
		logger.Sync()
		libvirtInst.LibvirtInst.Close()
	}()
	select {}
}
