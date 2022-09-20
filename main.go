package main

import (
	"IpProjectGo/ip"
	"IpProjectGo/logger"
	"IpProjectGo/ratelimit"

	"github.com/gin-gonic/gin"
)

func main() {
	ipServer := gin.Default()
	logServer := gin.Default()
	rateLimitServer := gin.Default()
	ipServer.POST("/", ip.Ip)
	logServer.POST("/", logger.Log)
	rateLimitServer.POST("/", ratelimit.RateLimitIp)
	go ipServer.Run("127.0.0.1:8000")
	go logServer.Run("127.0.0.1:6000")
	rateLimitServer.Run("127.0.0.1:5000")
}
