package main

import (
	"client/util"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

//"github.com/streadway/amqp"

var clients int = 1
var message string = "Hello World!\nHow are You?"

const (
	reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

const (
	count = 100 //10_000
)

func clientGO(conn util.Connection, topic_name string, wg *sync.WaitGroup, start <-chan struct{}) {
	defer conn.Disconnect()
	defer wg.Done()
	//msgBytes := util.RequestToJson(util.Request{Content: message })

	<-start

	for i := 0; i < count; i++ {

		msg := fmt.Sprintf("message %d", i+1)
		msgBytes := util.RequestToJson(util.Request{Content: msg, ResponseTo: topic_name})
		conn.Publish(topic_name, msgBytes)
		response := <-conn.Message
		log.Println("recebido: " + string(response))
	}

}

func main() {
	handleArgs()

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < 1; i++ {
		wg.Add(1)
		conn := util.NewConnection(util.Url, "client_0")
		/*conn.MessageHandler = func(msg []byte) {
			fmt.Printf("< recebido <%s>.\n", string(msg))
		}*/

		topic_name := fmt.Sprintf("fila_%d", i+1)
		conn.Subscribe(topic_name)

		go clientGO(*conn, topic_name, &wg, start)

		/*msg := ("Hello World <" + strconv.Itoa(i) + ">.")
		fmt.Printf("> enviado <%s>.\n", msg)
		conn.Publish("topico", []byte(msg))*/
	}

	time.Sleep(1 * time.Second)
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
