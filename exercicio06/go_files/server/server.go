package main

import (
	//connection "main/connection_mqqt"
	connection "main/connection_rabbitmq"
	util "main/util"
	"strings"
)

func countLines(str string) int {
	return 1 + int(strings.Count(str, "\n"))
}

func handleConnection(conn *connection.Connection, msg []byte) {
	request := util.JsonToRequest(msg)
	response := util.ResponseToJson(
		util.Response{
			Lines: countLines(request.Content),
		})
	conn.Publish(request.ResponseTo, response)
}

func main() {
	conn := connection.NewConnection("servidor")
	defer conn.Disconnect()
	//if len(os.Args) > 1 && os.Args[1] == "create" {
	conn.CreateQueue(util.Queue)
	//}
	conn.Subscribe(util.Queue)
	for {
		msg := <-conn.Message
		go handleConnection(conn, msg)
	}
}
