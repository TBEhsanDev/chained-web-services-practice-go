package test

import (
	"IpProjectGo/ip"
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

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TestUrl string = "http://127.0.0.1/api/"

const RequestNumber = 200

func TestWithJson(t *testing.T) {
	logFileLinesNumBefore := LogFileLines()
	mongodbLinesNumBefore := MongodbLInesCount()
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
	time.Sleep(4 * time.Second)
	assert.Equal(t, MongodbLInesCount(), mongodbLinesNumBefore+RequestNumber)

}

func TestWithoutJson(t *testing.T) {
	logFileLinesNumBefore := LogFileLines()
	mongodbLinesNumBefore := MongodbLInesCount()

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
	time.Sleep(4 * time.Second) //for logstash
	assert.Equal(t, MongodbLInesCount(), mongodbLinesNumBefore+RequestNumber)

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
	/*cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
		defer cursor.Close(context.TODO())
	} else {
		for cursor.Next(context.TODO()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				fmt.Println("cursor.Next() error:", err)
				os.Exit(1)
			} else {
			}
		}
	}*/
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
