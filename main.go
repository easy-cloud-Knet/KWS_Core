package main

import (
	_ "log"

	"github.com/easy-cloud-Knet/KWS_Core.git/api"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/router"
)

func main() {

	go api.Server(8080)
	router.MakeNewConnect()
}
