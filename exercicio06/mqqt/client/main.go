package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	//"time"
	connection "client/connection"
	util "client/util"
)

//"github.com/streadway/amqp"

var clients int = 1
var message string = "Hello World!\nHow are You?"

const (
	count = 10_000 //10_000
)

func clientGO(conn connection.Connection, responseTo string, wg *sync.WaitGroup, start <-chan struct{}) {
	defer conn.Disconnect()
	defer wg.Done()

	request := util.Request{
		Content:    message,
		ResponseTo: responseTo,
	}
	msgBytes := util.RequestToJson(request)

	<-start
	for i := 0; i < count; i++ {
		start := time.Now()
		///////////
		conn.Publish(util.Queue, msgBytes)
		<-conn.Message
		///////////
		delta := time.Since(start) / time.Nanosecond
		fmt.Println(strconv.FormatInt(delta.Nanoseconds(), 10))
		//log.Printf("%s[%s] recebido <%s>.%s", util.Yellow, responseTo, string(response), util.Reset)
	}
	//log.Printf("%s[%s] acabou.%s", util.Blue, responseTo, util.Reset)
}

func main() {
	handleArgs()

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < clients; i++ {
		wg.Add(1)
		conn := connection.NewConnection(
			util.Url,
			fmt.Sprintf("client_%d", i+1),
		)

		topic_name := fmt.Sprintf("fila_%d", i+1)
		conn.Subscribe(topic_name)

		go clientGO(*conn, topic_name, &wg, start)
	}

	//time.Sleep(1 * time.Second)
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

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
