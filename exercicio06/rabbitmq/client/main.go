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

	"github.com/streadway/amqp"
)

var clients int = 1
var message string = "Hello World!\nHow are You?"

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

const (
	count = 10_000
)

type Request struct {
	Content string
}

func connect(url string) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial(url)
	handleError(err, "🟥 conexão: %v")
	log.Printf("%s✅ conectado!%s", Blue, Reset)
	ch, err := conn.Channel()
	handleError(err, "🟥 canal: %v")
	return ch, conn
}

func toJson(request Request) []byte {
	msgBytes, err := json.Marshal(request)
	handleError(err, "🟥 marshal: %v")
	return msgBytes
}

func clientGO(wg *sync.WaitGroup, start <-chan struct{}, id int) {
	defer wg.Done()
	<-start

	ch, conn := connect("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	defer ch.Close()

	msgBytes := toJson(Request{Content: message})

	err := ch.Publish(
		"",
		"exercicio06",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgBytes,
		},
	)

	handleError(err, "🟥 Publicar: %v")
	log.Printf("🟩 mensagem enviada %s%s%s.", Yellow, msgBytes, Reset)
}

func main() {
	handleArgs()

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < clients; i++ {
		wg.Add(1)
		go clientGO(&wg, start, i+1)
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
