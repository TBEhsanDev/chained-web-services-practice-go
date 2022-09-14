package main

import (
	"IpProjectGo/ip"

	"github.com/gin-gonic/gin"
)

func main() {
	ipServer := gin.Default()
	ipServer.POST("/", ip.Ip)
	ipServer.Run("127.0.0.1:8000")
}
