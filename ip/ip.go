package ip

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type D struct {
	Student string `json`
	Ip      string `json`
}

func Ip(c *gin.Context) {
	var data D
	Ip := c.Request.Header.Get("X-Real-IP")
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
	}
	data.Ip = Ip
	c.JSON(200, data)
}
