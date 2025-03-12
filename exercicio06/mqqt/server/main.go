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

func handleConnection(ch *amqp.Channel, msg amqp.Delivery) {
	request := util.JsonToRequest(msg.Body)
	log.Printf("> mensagem %srecebida %s%s.", yellow, reset, strings.Replace(string(msg.Body), "\n", "\\n", -1))
	response := util.ResponseToJson(
		util.Response{
			Lines: countLines(request.Content),
		})
	util.Publish(ch, msg.ReplyTo, response)
	log.Printf("< mensagem  %senviada %s%s.", green, reset, response)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "create" {
		util.CreateQueue("exercicio06")
	}

	conn := util.NewConnection(util.Url)
	defer conn.Close()
	ch := util.NewChannel(conn)
	defer ch.Close()
	msgs := util.NewConsumer(ch, util.Queue)

	for {
		msg := <-msgs
		go handleConnection(ch, msg)
	}
}
