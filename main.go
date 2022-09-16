package main

import (
	"IpProjectGo/ip"
	"IpProjectGo/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	ipServer := gin.Default()
	ipServer.POST("/", ip.Ip)
	logServer := gin.Default()
	logServer.POST("/", logger.Log)
	go ipServer.Run("127.0.0.1:8000")
	logServer.Run("127.0.0.1:6000")
}
