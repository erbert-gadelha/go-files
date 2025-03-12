package main

import (
	"log"
	"os"
	"strings"

	util "server/util"

	"github.com/streadway/amqp"
)

const (
	reset   = "\033[0m"
	Red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

func countLines(str string) int {
	return 1 + int(strings.Count(str, "\n"))
}

func newConnection(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	handleError(err, "🟥 conexão: %v")
	log.Printf("%s✅ conectado!%s", Blue, reset)
	return conn
}

func newChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	handleError(err, "🟥 canal: %v")
	return ch
}

func newConsumer(ch *amqp.Channel) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		"exercicio06",
		"", true, false,
		false, false, nil,
	)
	handleError(err, "🟥 consumidor: %v")
	return msgs
}

func publish(ch *amqp.Channel, msg []byte) {
	err := ch.Publish(
		"",
		"exercicio06",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
	handleError(err, "🟥 Publicar: %v")
}

func handleConnection() {
	conn := util.NewConnection(util.Url)
	defer conn.Close()
	ch := util.NewChannel(conn)
	defer ch.Close()
	msgs := util.NewConsumer(ch, util.Queue)

	for {
		msg := <-msgs
		request := util.JsonToRequest(msg.Body)
		log.Printf("> mensagem %srecebida %s%s.", yellow, reset, strings.Replace(string(msg.Body), "\n", "\\n", -1))
		response := util.ResponseToJson(
			util.Response{
				Lines: countLines(request.Content),
			})

		util.Publish(ch, msg.ReplyTo, response)
		log.Printf("< mensagem  %senviada %s%s.", green, reset, response)

	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "create" {
		createQueue(
			"exercicio06", // name
			false,         // durable
			false,         // autoDelete
			false,         // exclusive
			false,         // noWait
			nil,           // args
		)
	}
	handleConnection()
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}

func createQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("| fila <%s> ⭕ %v", name, err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("| fila <%s> ⭕ %v", name, err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		name, durable, autoDelete, exclusive, noWait, args,
	)

	if err != nil {
		log.Fatalf("| fila <%s> ⭕ %v", name, err)
	}
	log.Printf("✅ %scriada fila <%s>%s", Blue, name, reset)
}
