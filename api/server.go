package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Server(portNum int){
	router:=gin.Default()

	router.GET("/gin", func(c *gin.Context){
		fmt.Println(c)
	})

	router.GET("/gin", func(c *gin.Context){
		fmt.Println(c)
	})
	router.GET("/gin", func(c *gin.Context){
		fmt.Println(c)
	})
	router.GET("/gin", func(c *gin.Context){
		fmt.Println(c)
	})
	router.GET("/gin", func(c *gin.Context){
		fmt.Println(c)
	})
	router.GET("/gin", func(c *gin.Context){
		fmt.Println(c)
	})
	router.Run(":"+strconv.Itoa(portNum))
}