package main

import (
	"fmt"
	"io"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"
)

var clients int = 1
var message string = "Hello World!\nHow are You?"

const (
	count = 10_000
)

type Request struct {
	Content string
}

func callRemote(message string, client *rpc.Client) {
	request := Request{Content: message}
	var response int
	err := client.Call("Arquivo.CountLines", request, &response)
	if err != nil {
		print(err)
	}
}

func clientGO(client *rpc.Client, wg *sync.WaitGroup, start <-chan struct{}, id int) {
	defer wg.Done()
	defer client.Close()
	<-start

	for i := 0; i < count; i++ {
		start := time.Now()
		callRemote(message, client)
		delta := time.Since(start) / time.Nanosecond
		fmt.Println(strconv.FormatInt(delta.Nanoseconds(), 10))
	}

}

func main() {

	if len(os.Args) > 1 {
		clients, _ = strconv.Atoi(os.Args[1])
	}

	if len(os.Args) > 2 {
		message = readFile(os.Args[2])
	}

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < clients; i++ {
		client, err := rpc.Dial("tcp", "localhost:1313")
		if err != nil {
			print(err)
		}

		wg.Add(1)
		go clientGO(client, &wg, start, i+1)
	}

	time.Sleep(1 * time.Second)
	close(start)

	wg.Wait()
	fmt.Println()
}

func readFile(fileName string) string {
	file, _ := os.Open(fileName)
	content, _ := io.ReadAll(file)
	file.Close()
	return string(content)
}
