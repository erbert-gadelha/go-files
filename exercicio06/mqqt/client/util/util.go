package util

import (
	"encoding/json"
	"log"
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
	Content    string
	ResponseTo string
}
type Response struct {
	Lines int
}

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
