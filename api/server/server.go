package server

import (
	"fmt"
	"strconv"

	"libvirt.org/go/libvirt"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"github.com/gin-gonic/gin"
)

func InitServer(portNum int, libvirtInst *libvirt.Connect){
	router:=gin.Default()

	
	router.GET("/getStatus", func(c *gin.Context){
		fmt.Println("getStatus request income")
		conn.ActiveDomain(libvirtInst)
	})
	router.Run(":"+strconv.Itoa(portNum))
}
