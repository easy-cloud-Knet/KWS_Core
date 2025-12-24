package main

import (
	_ "context"
	"fmt"
	_ "log"
	"runtime/debug"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core/api"
	syslogger "github.com/easy-cloud-Knet/KWS_Core/logger"
	"github.com/easy-cloud-Knet/KWS_Core/server"
	snapmgrpkg "github.com/easy-cloud-Knet/KWS_Core/vm/service/snapshot"
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
	libvirtInst.DomainControl.DomainListStatus.UpdateCPUTotal()
	libvirtInst.DomainControl.RetrieveAllDomain(libvirtInst.LibvirtInst, logger)

	// Inject SnapshotManager into the InstHandler
	{
		snapmgr := snapmgrpkg.NewManagerWithDeps(domListCon, libvirtInst.LibvirtInst, "/var/lib/kws/snapshots")
		libvirtInst.SnapshotManager = snapmgr
	}

	go server.InitServer(8080, &libvirtInst, *logger)
	fmt.Println("asfd")
	defer func() {
		logger.Info("Shutting down gracefully...") // 종료 시 로깅
		logger.Sync()
		libvirtInst.LibvirtInst.Close()
	}()
	select {}
}
