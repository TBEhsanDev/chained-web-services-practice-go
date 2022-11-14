package logger

import (
	"IpProjectGo/ip"
	"IpProjectGo/ratelimit"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger struct {
	Time string
	Data string
}

func Log(c *gin.Context) {
	var error500 ratelimit.Error
	error500.Detail = "Bad Gateway"
	var logData Logger
	var data ip.D
	GetRequestData(c, &data, &logData)
	writeLog(logData)
	response, err := RequestToIp(data)
	if err == nil {
		defer response.Body.Close()
		respData := ReturnResponse(response)
		c.JSON(200, respData)
	} else {
		c.JSON(500, error500)
	}
}
func GetRequestData(c *gin.Context, data *ip.D, logData *Logger) {
	Ip := c.Request.Header.Get("X-Real-IP")
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
	}
	data.Ip = Ip
	logData.Data = data.Data
	logData.Time = time.Now().String()
}
func ReturnResponse(response *http.Response) ip.D {

	respJson, _ := ioutil.ReadAll(response.Body)
	respData := ip.D{}
	if err := json.Unmarshal(respJson, &respData); err != nil {
		fmt.Println(err)
	}
	return respData
}
func RequestToIp(data ip.D) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	request, err := http.NewRequest("POST", "http://127.0.0.1:8000", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("X-Real-IP", data.Ip)
	client := &http.Client{}
	response, err := client.Do(request)
	return response, err
}
func writeLog(logData Logger) {
	f, err := os.OpenFile("./log.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.Encode(logData)
}
