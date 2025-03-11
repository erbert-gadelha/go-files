package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
)

type Request struct {
	Content string
}

type Arquivo struct {
	logging bool
}

func (a *Arquivo) CountLines(request *Request, response *int) error {
	*response = (1 + int(strings.Count(request.Content, "\n")))
	return nil
}

func listenRPC(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()

	fmt.Printf("Servidor RPC rodando em (%s)...\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func main() {
	logging := false
	if len(os.Args) > 1 && os.Args[1] == "--log" {
		logging = true
	}
	arquivo := &Arquivo{logging: logging}
	rpc.Register(arquivo)
	listenRPC("0.0.0.0:1313")
}
