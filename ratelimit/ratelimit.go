package ratelimit

import (
	"IpProjectGo/ip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Error struct {
	Detail string
}

func RateLimitIp(c *gin.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var error429 Error
	var error500 Error
	error500.Detail = "Bad Gateway"
	error429.Detail = "too many requests"
	var m sync.Mutex
	m.Lock()
	check := checkRate(c, client)
	m.Unlock()
	if check == false {
		c.JSON(429, error429)
	} else {
		var data ip.D
		GetData(c, &data)
		response, err := RequestToLoggerApp(data)
		if err != nil || response.StatusCode == 500 {
			c.JSON(500, error500)
		} else {
			defer response.Body.Close()
			respData := ReturnSentResponse(response)
			c.JSON(200, respData)
		}
	}
}
func checkRate(c *gin.Context, client *redis.Client) bool {
	Ip := c.Request.Header.Get("X-Real-IP")
	ip := string(Ip)
	v, err := client.Get(ip).Result()
	if err != nil {
		err = client.Set(ip, 100, 1*time.Minute).Err()
		if err != nil {
			fmt.Println("1")
		}
	}
	num, err := strconv.Atoi(v)
	if num > 0 {
		client.Decr(ip)
	}
	print(num)
	if num == 0 {
		return false
	}
	return true
}
func GetData(c *gin.Context, data *ip.D) {

	Ip := c.Request.Header.Get("X-Real-IP")
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
	}
	data.Ip = Ip
}
func ReturnSentResponse(response *http.Response) ip.D {
	respJson, _ := ioutil.ReadAll(response.Body)
	respData := ip.D{}
	if err := json.Unmarshal(respJson, &respData); err != nil {
		fmt.Println(err)
	}
	return respData
}
func RequestToLoggerApp(data ip.D) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	request, err := http.NewRequest("POST", "http://127.0.0.1:6000", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("X-Real-IP", data.Ip)
	client := &http.Client{}
	response, err := client.Do(request)
	fmt.Println(response.StatusCode)
	return response, err
}
