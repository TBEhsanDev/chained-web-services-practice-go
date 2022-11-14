package test

import (
	"IpProjectGo/ip"
	"IpProjectGo/ratelimit"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TestUrl string = "http://127.0.0.1/api/"

const RequestNumber = 200
const AllowedRequestsNumber = 100

var checkStatus = true

func TestWithJson(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	logFileLinesNumBefore := LogFileLines()
	mongodbLinesNumBefore := MongodbLInesCount()
	reqJson, respJson := Setup(ip.D{Data: "ali"})
	client.Del("127.0.0.1")
	for i := 0; i < RequestNumber; i++ {
		var response *http.Response
		myChannel := make(chan *http.Response)
		go Request(t, reqJson, respJson, myChannel)
		response = <-myChannel
		defer response.Body.Close()
	}
	client.Del("127.0.0.1")
	if checkStatus == true {
		assert.Equal(t, logFileLinesNumBefore+AllowedRequestsNumber, LogFileLines())
		time.Sleep(5 * time.Second) //for logstash
		assert.Equal(t, mongodbLinesNumBefore+AllowedRequestsNumber, MongodbLInesCount())
	}
}
func TestWithoutJson(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	logFileLinesNumBefore := LogFileLines()
	mongodbLinesNumBefore := MongodbLInesCount()
	reqJson, respJson := Setup(ip.D{})
	client.Del("127.0.0.1")
	for i := 0; i < RequestNumber; i++ {
		var response *http.Response
		myChannel := make(chan *http.Response)
		go Request(t, reqJson, respJson, myChannel)
		response = <-myChannel
		defer response.Body.Close()
	}
	client.Del("127.0.0.1")
	if checkStatus == true {
		assert.Equal(t, logFileLinesNumBefore+AllowedRequestsNumber, LogFileLines())
		time.Sleep(5 * time.Second) //for logstash
		assert.Equal(t, mongodbLinesNumBefore+AllowedRequestsNumber, MongodbLInesCount())
	}
}
func MongodbLInesCount() int64 {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}
	coll := client.Database("ip_db").Collection("ip_collection")
	estCount, estCountErr := coll.EstimatedDocumentCount(context.TODO())
	if estCountErr != nil {
		panic(estCountErr)
	}
	return estCount
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
func Request(t *testing.T, reqJson []byte, respJson []byte, myChannel chan *http.Response) {
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
	if response.StatusCode == 500 {
		checkStatus = false
		var e ratelimit.Error
		respJson, _ = ioutil.ReadAll(response.Body)
		if err := json.Unmarshal(respJson, &e); err != nil {
			fmt.Println(err)
		}
		assert.Equal(t, "Bad Gateway", e.Detail)
	}
	if response.StatusCode == 429 {
		var e ratelimit.Error
		respJson, _ = ioutil.ReadAll(response.Body)
		if err := json.Unmarshal(respJson, &e); err != nil {
			fmt.Println(err)
		}
		assert.Equal(t, "too many requests", e.Detail)
	}
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		assert.Equal(t, body, respJson)
	}
	myChannel <- response
}
