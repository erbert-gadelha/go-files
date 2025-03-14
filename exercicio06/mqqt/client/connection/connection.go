package connection

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Connection struct {
	client         mqtt.Client
	Message        chan []byte
	MessageHandler func([]byte)
}

func (c *Connection) Publish(queue string, msg []byte) {
	//token := c.client.Publish(queue, 1, false, msg)
	token := c.client.Publish(queue, 0, false, msg)
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
	//fmt.Printf("inscrito no tópico <%s>\n", queue)
}

func (c *Connection) Disconnect() {
	c.client.Disconnect(255)
}

func NewConnection(url string, id string) *Connection {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(url)
	opts.SetClientID(id)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	conn := Connection{
		client:         client,
		Message:        make(chan []byte),
		MessageHandler: nil,
	}

	conn.MessageHandler = func(b []byte) {
		conn.Message <- b
	}

	opts.SetDefaultPublishHandler(func(mqtt_client mqtt.Client, msg mqtt.Message) {
		conn.MessageHandler(msg.Payload())
	})
	return &conn
}
