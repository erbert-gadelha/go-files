package main

import (
	"os"
	connection "server/connection"
	util "server/util"
	"strings"
)

func countLines(str string) int {
	return 1 + int(strings.Count(str, "\n"))
}

func handleConnection(conn *connection.Connection, msg []byte) {
	request := util.JsonToRequest(msg)
	//content := strings.Replace(string(request.Content), "\n", "\\n", -1)
	//log.Printf("%s[%s] > %s%s.", util.Yellow, request.ResponseTo, content, util.Reset)
	response := util.ResponseToJson(
		util.Response{
			Lines: countLines(request.Content),
		})

	conn.Publish(request.ResponseTo, response)
	//log.Printf("%s[%s] < %s%s", util.Green, request.ResponseTo, string(response), util.Reset)
}

func main() {

	conn := connection.NewConnection(util.Url, "servidor")
	defer conn.Disconnect()

	if len(os.Args) > 1 && os.Args[1] == "create" {
		conn.CreateQueue(util.Queue)
	}

	conn.Subscribe(util.Queue)
	for {
		msg := <-conn.Message
		go handleConnection(conn, msg)
	}
}
