package util

import (
	"encoding/json"
	"log"
)

const (
	Reset        = "\033[0m"
	Red          = "\033[31m"
	Green        = "\033[32m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Magenta      = "\033[35m"
	Cyan         = "\033[36m"
	Gray         = "\033[37m"
	White        = "\033[97m"
	URI_rabbitMQ = "amqp://guest:guest@localhost:5672/"
	URI_mqqt     = "tcp://localhost:1883"
	Queue        = "exercicio06"
)

type Request struct {
	Content    string
	ResponseTo string
}
type Response struct {
	Lines int
}

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf(Red+msg+Reset, err)
	}
}

func ResponseToJson(request Response) []byte {
	msgBytes, err := json.Marshal(request)
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
