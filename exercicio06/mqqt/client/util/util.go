package util

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	gray    = "\033[37m"
	white   = "\033[97m"
	Url     = "tcp://localhost:1883"
	Queue   = "exercicio06"
)

type Request struct {
	Content string
}
type Response struct {
	Lines int
}

type Connection struct {
	client         mqtt.Client
	MessageHandler func([]byte)
}

func (c *Connection) Publish(queue string, msg []byte) {
	token := c.client.Publish(queue, 1, false, msg)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
func (c *Connection) Subscribe(queue string) {
	if token := c.client.Subscribe(queue, 1, func(client mqtt.Client, msg mqtt.Message) {
		c.MessageHandler(msg.Payload())
	}); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("inscrito no tópico <%s>\n", queue)
}

func (c *Connection) Disconnect() {
	c.client.Disconnect(255)
}

func NewConnection(url string, id string) *Connection {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(Url)
	opts.SetClientID(id)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	conn := Connection{client: client, MessageHandler: nil}

	opts.SetDefaultPublishHandler(func(mqtt_client mqtt.Client, msg mqtt.Message) {
		conn.MessageHandler(msg.Payload())
	})
	return &conn
}

func Funcao() {
	conn := NewConnection(Url, "client_0")
	conn.MessageHandler = func(msg []byte) {
		fmt.Printf("< recebido <%s>.\n", string(msg))
	}

	conn.Subscribe("topico")

	for i := 0; i < 10; i++ {
		msg := fmt.Sprint("Hello <%d>.", i)
		fmt.Printf("> enviado <%s>.\n", msg)
		conn.Publish("topico", []byte(msg))
		time.Sleep(1 * time.Second)
	}
	conn.Disconnect()
	fmt.Println("<Fin>")
}

/*
func NewConnection(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	HandleError(err, "🟥 conexão: %v")
	log.Printf("%s✅ conectado!%s", blue, reset)
	return conn
}

func NewChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	HandleError(err, "🟥 canal: %v")
	return ch
}

func NewConsumer(ch *amqp.Channel, queue string) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		queue,
		"", true, false,
		false, false, nil,
	)
	HandleError(err, "🟥 consumidor: %v")
	return msgs
}

func Publish(ch *amqp.Channel, replyTo string, msg []byte) {
	err := ch.Publish(
		"",
		"exercicio06",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
			ReplyTo:     replyTo,
		},
	)
	HandleError(err, "🟥 Publicar: %v")
}

func HandleConnection(url string, queue string) {
	conn := NewConnection(url)
	defer conn.Close()
	ch := NewChannel(conn)
	defer ch.Close()
	msgs := NewConsumer(ch, queue)

	for {
		msgBytes := <-msgs
		log.Printf(" mensagem recebida (%d) %s%s%s.", 2, yellow, string(msgBytes.Body), reset)
	}
}



func CreateQueue(name string) {
	CreateQueue_(name, false, false, false, false, nil)
}

func CreateQueue_(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) {
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
	log.Printf("✅ %scriada fila <%s>%s", blue, name, reset)
}*/

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
func ResponseToJson(response Response) []byte {
	msgBytes, err := json.Marshal(response)
	HandleError(err, "🟥 marshal: %v")
	return msgBytes
}

func RequestToJson(request Request) []byte {
	msgBytes, err := json.Marshal(request)
	HandleError(err, "🟥 marshal: %v")
	return msgBytes
}

func JsonToResponse(response []byte) Response {
	r := Response{}
	err := json.Unmarshal(response, &r)
	HandleError(err, "🟥 unmarshal: %v")
	return r
}

func JsonToRequest(request []byte) Request {
	r := Request{}
	err := json.Unmarshal(request, &r)
	HandleError(err, "🟥 unmarshal: %v")
	return r
}
