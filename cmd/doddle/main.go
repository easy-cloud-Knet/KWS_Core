package main

import (
	_ "context"
	_ "log"
	"runtime/debug"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core/api"
	create "github.com/easy-cloud-Knet/KWS_Core/api/Create"
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

	// TODO: 모든 Handler의 생성자에서 DomainControl과 LibvirtConnect를 주입받도록 수정
	createHandler := create.NewHandler(domListCon, libvirtInst.LibvirtInst, logger)

	libvirtInst.LibvirtConnection()
	libvirtInst.DomainControl.SetLibvirt(libvirtInst.LibvirtInst)
	libvirtInst.DomainControl.DomainListStatus.Update()
	if err := libvirtInst.DomainControl.RetrieveAllDomain(logger); err != nil {
		logger.Fatal("failed to retrieve domains on startup", zap.Error(err))
	}

	go server.InitServer(config.ServerPort, &libvirtInst, createHandler, logger)
	defer func() {
		logger.Info("Shutting down gracefully...") // 종료 시 로깅
		logger.Sync()
		libvirtInst.LibvirtInst.Close()
	}()
	select {}
}
