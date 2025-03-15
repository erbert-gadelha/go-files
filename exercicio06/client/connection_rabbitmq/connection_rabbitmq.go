package connection_rabbitmq

import (
	util "client/util"
	"log"

	amqp "github.com/streadway/amqp"
)

type Connection struct {
	conn           *amqp.Connection
	ch             *amqp.Channel
	replyTo        string
	Message        chan []byte
	MessageHandler func([]byte)
}

func (c *Connection) Publish(queue string, msg []byte) {
	err := c.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
			ReplyTo:     c.replyTo,
		},
	)
	util.HandleError(err, "🟥 [RABBIT] Publicar: %v")
}

func (c *Connection) Subscribe(queue string) {
	msgs := newConsumer(c.ch, queue)
	go (func() {
		for {
			msg := <-msgs
			c.Message <- msg.Body
		}
	})()
}

func (c *Connection) Disconnect() {
	c.ch.Close()
	c.conn.Close()
}

func NewConnection(url string, id string) *Connection {
	conn := newConn(url)
	c := Connection{
		conn:           conn,
		ch:             newChannel(conn),
		replyTo:        id,
		Message:        make(chan []byte),
		MessageHandler: nil,
	}
	return &c
}

func newConn(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	util.HandleError(err, "🟥 [RABBIT] conexão: %v")
	log.Printf("%s✅ [RABBIT] conectado!%s", util.Blue, util.Reset)
	return conn
}

func newChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	util.HandleError(err, "🟥 [RABBIT] canal: %v")
	return ch
}

func newConsumer(ch *amqp.Channel, queue string) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		queue,
		"", true, false,
		false, false, nil,
	)
	util.HandleError(err, "🟥 [RABBIT] consumidor: %v")
	return msgs
}

func (c *Connection) CreateQueue(name string) {
	CreateQueue(name, false, false, false, false, nil, c)
}

func CreateQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table, connection *Connection) {
	_, err := connection.ch.QueueDeclare(
		name, durable, autoDelete, exclusive, noWait, args,
	)
	util.HandleError(err, "⭕ [RABBIT] criar fila> %v")
	log.Printf("✅%s [RABBIT] criada fila <%s>%s", util.Blue, name, util.Reset)
}
