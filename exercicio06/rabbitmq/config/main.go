package main

import (
	"log"

	"github.com/streadway/amqp"
)

func configQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) {
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
	log.Printf("| fila <%s> ✅", name)
}

func main() {
	configQueue(
		"exercicio06", // name
		false,         // durable
		false,         // autoDelete
		false,         // exclusive
		false,         // noWait
		nil,           // args
	)
}
