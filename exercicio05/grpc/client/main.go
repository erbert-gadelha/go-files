package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	pb "client.go/grpcarquivo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clients int = 1
var message string = "Hello World!\nHow are You?"
var ctx, cancel = context.WithCancel(context.Background())

const (
	count = 10_000
)

func callRemote(message string, client pb.ArquivoServiceClient) {
	req := &pb.Request{Content: message}
	_, err := client.CountLines(ctx, req)
	if err != nil {
		print(err)
	}
}

func clientGO(client pb.ArquivoServiceClient, wg *sync.WaitGroup, start <-chan struct{}, id int) {
	defer wg.Done()
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
		conn, err := grpc.NewClient("localhost:1313", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			print(err)
		}

		defer conn.Close()
		client := pb.NewArquivoServiceClient(conn)

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
