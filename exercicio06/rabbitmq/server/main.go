package main

import (
	"log"
	"net/rpc"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

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

type Request struct {
	Content string
}

type Arquivo struct {
	logging bool
}

func countLines(str string) int {
	return 1 + int(strings.Count(str, "\n"))
}

func connect(url string) (*amqp.Channel, *amqp.Connection, <-chan amqp.Delivery) {
	conn, err := amqp.Dial(url)
	handleError(err, "🟥 conexão: %v")
	log.Printf("%s✅ conectado!%s", Blue, Reset)
	ch, err := conn.Channel()
	handleError(err, "🟥 canal: %v")
	msgs, err := ch.Consume(
		"exercicio06",
		"", true, false,
		false, false, nil,
	)
	handleError(err, "🟥 conumidor: %v")
	return ch, conn, msgs
}

func handleConnection(url string) {
	ch, conn, msgs := connect(url)
	defer ch.Close()
	defer conn.Close()

	for {
		msgBytes := <-msgs
		log.Printf(" mensagem recebida %s%s%s.", Yellow, string(msgBytes.Body), Reset)

	}
}

func main() {
	logging := false
	if len(os.Args) > 1 && os.Args[1] == "--log" {
		logging = true
	}
	arquivo := &Arquivo{logging: logging}
	rpc.Register(arquivo)
	handleConnection("amqp://guest:guest@localhost:5672/")
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
