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
		app.InstHandler,
		app.ControlHandler,
		app.CreateHandler,
		app.SnapshotHandler,
		app.Logger)

	select {}
}
