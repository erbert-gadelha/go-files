package main

import (
	"fmt"
	"io"
	"os"
	"log"
	"net/rpc"
)

type Args struct {
	Conteudo string
}


func readFile(fileName string) string {
	file, _ := os.Open(fileName)
	content, _ := io.ReadAll(file)
	file.Close()
	return string(content)
}

func sendText(content string) *int {
	client, err := rpc.Dial("tcp", "localhost:1313")
	if err != nil { log.Fatalf("Erro ao conectar: %v", err)	}
	defer client.Close()

	args := Args{Conteudo:content}
	var reply int

	err = client.Call("Arquivo.Linhas", args, &reply)
	if err != nil { log.Fatalf("Erro ao chamar RPC: %v", err) }

	return &reply
}



func main() {
	content:="hello world!\nHow you doing?\n"
	if len(os.Args) > 1 {
		content = readFile(os.Args[1])
	}

	reply := sendText(content)
	fmt.Printf("resultado %d.\n", *reply)
}
