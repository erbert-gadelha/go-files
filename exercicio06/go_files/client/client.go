package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	util "main/util"

	connection "main/connection_rabbitmq"
	//connection "main/connection_mqqt"
)

var clients int = 1
var message string = "Hello World!\nHow are You?"

const (
	count = 10_000 //10_000
)

func clientGO(conn connection.Connection, responseTo string, wg *sync.WaitGroup, start <-chan struct{}) {
	defer conn.Disconnect()
	defer wg.Done()
	msgBytes := util.RequestToJson(
		util.Request{
			Content:    message,
			ResponseTo: responseTo,
		},
	)
	<-start
	for i := 0; i < count; i++ {
		start := time.Now()
		///////////
		conn.Publish(util.Queue, msgBytes)
		<-conn.Message
		///////////
		delta := time.Since(start) / time.Nanosecond
		fmt.Println(strconv.FormatInt(delta.Nanoseconds(), 10))
	}
}

func main() {
	handleArgs()
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < clients; i++ {
		wg.Add(1)
		id := fmt.Sprintf("client_%d", i+1)
		conn := connection.NewConnection(id)
		topic_name := fmt.Sprintf("fila_%d", i+1)
		conn.CreateQueue(topic_name)
		conn.Subscribe(topic_name)
		go clientGO(*conn, topic_name, &wg, start)
	}
	close(start)
	wg.Wait()
	fmt.Println()
}

func handleArgs() {
	if len(os.Args) > 1 {
		clients, _ = strconv.Atoi(os.Args[1])
	}

	if len(os.Args) > 2 {
		message = readFile(os.Args[2])
	}
}

func readFile(fileName string) string {
	file, _ := os.Open(fileName)
	content, _ := io.ReadAll(file)
	file.Close()
	return string(content)
}
