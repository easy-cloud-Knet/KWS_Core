package main

import (
	"runtime/debug"

	"github.com/easy-cloud-Knet/KWS_Core/internal/config"
	"github.com/easy-cloud-Knet/KWS_Core/internal/server"
)

func main() {
	debug.SetTraceback("none")

	app := initApp()
	defer app.Shutdown()

	go server.InitServer(
		config.ServerPort,
		server.Handlers{
			IsConnected: app.IsConnected,
			Control:     app.ControlHandler,
			Create:      app.CreateHandler,
			Snapshot:    app.SnapshotHandler,
			Status:      app.StatusHandler,
		},
		app.Logger)

	select {}
}
