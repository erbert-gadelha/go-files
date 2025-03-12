package main

import (
	"log"
	"os"
	"strings"
	util "server/util"
	connection "server/connection"
)


func countLines(str string) int {
	return 1 + int(strings.Count(str, "\n"))
}

func handleConnection(conn *connection.Connection, msg []byte) {
	request := util.JsonToRequest(msg)
	log.Printf("> mensagem %srecebida %s%s.", util.Yellow, util.Reset, strings.Replace(string(msg), "\n", "\\n", -1))
	response := util.ResponseToJson(
		util.Response{
			Lines: countLines(request.Content),
		})

	conn.Publish(request.ResponseTo, response)
	log.Printf("< mensagem  %senviada %s%s.", util.Green, util.Reset, response)
}

func main() {

	conn := connection.NewConnection(util.Url, util.Queue)
	defer conn.Disconnect()

	/*if len(os.Args) > 1 && os.Args[1] == "create" {
		conn.CreateQueue(util.Queue)
	}*/

	for {
		msg := <-conn.Message
		go handleConnection(conn, msg)
	}
}
