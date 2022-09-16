package test

import (
	"IpProjectGo/ip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var TestUrl string = "http://127.0.0.1/api/"

const RequestNumber = 200

func TestWithJson(t *testing.T) {
	logFileLinesNumBefore := LogFileLines()
	reqJson, respJson := Setup(ip.D{Student: "ali"})
	for i := 0; i < RequestNumber; i++ {
		var response *http.Response
		myChannel := make(chan *http.Response)
		go Request(reqJson, myChannel)
		response = <-myChannel
		defer response.Body.Close()
		CheckResponse(t, response, respJson)
	}
	assert.Equal(t, LogFileLines(), logFileLinesNumBefore+RequestNumber)
}

func TestWithoutJson(t *testing.T) {
	logFileLinesNumBefore := LogFileLines()
	reqJson, respJson := Setup(ip.D{})
	for i := 0; i < RequestNumber; i++ {
		var response *http.Response
		myChannel := make(chan *http.Response)
		go Request(reqJson, myChannel)
		response = <-myChannel
		defer response.Body.Close()
		CheckResponse(t, response, respJson)
	}
	assert.Equal(t, LogFileLines(), logFileLinesNumBefore+RequestNumber)

}
func LogFileLines() int {
	lineCount := 0
	file, _ := os.Open("../log.jsonl")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}
	return lineCount
}
func Setup(a ip.D) ([]byte, []byte) {
	var ReqData ip.D = a
	a.Ip = "127.0.0.1"
	var RespData ip.D = a
	respJson, err := json.Marshal(RespData)
	if err != nil {
		fmt.Println(err)
	}
	reqJson, err := json.Marshal(ReqData)
	if err != nil {
		fmt.Println(err)
	}
	return reqJson, respJson
}
func Request(reqJson []byte, myChannel chan *http.Response) {
	request, err := http.NewRequest("POST", TestUrl, bytes.NewBuffer(reqJson))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("HTTP call failed:", err)
	}
	myChannel <- response
}
func CheckResponse(t *testing.T, response *http.Response, respJson []byte) {
	if response.StatusCode == 500 {
		fmt.Println("Bad Gateway:", response.StatusCode)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, body, respJson)
}
