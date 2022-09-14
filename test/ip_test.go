package test

import (
	"IpProjectGo/ip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

var TestUrl string = "http://127.0.0.1/api/"

func TestWithJson(t *testing.T) {
	reqData := ip.D{Student: "ali"}
	respData := ip.D{Student: "ali", Ip: "127.0.0.1"}
	respJson, err := json.Marshal(respData)
	if err != nil {
		fmt.Println(err)
	}
	reqJson, err := json.Marshal(reqData)
	if err != nil {
		fmt.Println(err)
	}
	request, err := http.NewRequest("POST", TestUrl, bytes.NewBuffer(reqJson))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("HTTP call failed:", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Println("Bad Gateway:", response.StatusCode)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, body, respJson)
}

func TestWithoutJson(t *testing.T) {
	reqData := ip.D{}
	respData := ip.D{Ip: "127.0.0.1"}
	respJson, err := json.Marshal(respData)
	if err != nil {
		fmt.Println(err)
	}
	reqJson, err := json.Marshal(reqData)
	if err != nil {
		fmt.Println(err)
	}
	request, err := http.NewRequest("POST", TestUrl, bytes.NewBuffer(reqJson))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("HTTP call failed:", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Println("Bad Gateway:", response.StatusCode)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, body, respJson)
}
