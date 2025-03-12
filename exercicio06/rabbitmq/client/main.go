package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	util "client/util"

	"github.com/streadway/amqp"
)

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
	count = 10 //10_000
)

type Request struct {
	Content string
}

func connect(url string) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial(url)
	handleError(err, "🟥 conexão: %v")
	log.Printf("%s✅ conectado!%s", Blue, reset)
	ch, err := conn.Channel()
	handleError(err, "🟥 canal: %v")
	return ch, conn
}

func toJson(request Request) []byte {
	msgBytes, err := json.Marshal(request)
	handleError(err, "🟥 marshal: %v")
	return msgBytes
}

func clientGO(conn *amqp.Connection, ch *amqp.Channel, msgs <-chan amqp.Delivery, replyTo string, wg *sync.WaitGroup, start <-chan struct{}) {
	defer conn.Close()
	defer wg.Done()

	msgBytes := util.RequestToJson(util.Request{Content: message})

	for i := 0; i < count; i++ {
		util.Publish(ch, replyTo, msgBytes)
		msg := <-msgs

		response := util.JsonToResponse(msg.Body)

		log.Printf("🟩 linhas %s%d%s.", yellow, response.Lines, reset)
	}

}

func main() {
	handleArgs()

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < clients; i++ {
		wg.Add(1)
		queueName := fmt.Sprintf("fila_%d", i+1)

		conn := util.NewConnection(util.Url)
		ch := util.NewChannel(conn)
		util.CreateQueue(queueName)
		msgs := util.NewConsumer(ch, queueName)

		go clientGO(conn, ch, msgs, queueName, &wg, start)
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
