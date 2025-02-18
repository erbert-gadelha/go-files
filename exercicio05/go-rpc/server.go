package main

import (
	"fmt"
	"net/rpc"
	"net"
	"os"
)


type Args struct {
	Conteudo string
}

type Arquivo struct {
	logging bool
}

func (a *Arquivo) Linhas(args *Args, reply *int) error {
	*reply = len(args.Conteudo)
	if(a.logging) {
		fmt.Println(*args)
	}
	return nil
}

func listenRPC(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil { fmt.Println(err) }
	defer listener.Close()

	fmt.Printf("Servidor RPC rodando em (%s)...\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil { fmt.Println(err); continue }
		go rpc.ServeConn(conn)
	}
}

func main() {
	logging:=false
	if (len(os.Args)>1 && os.Args[1]=="--log") { logging=true }
	arquivo := &Arquivo{logging:logging}
	rpc.Register(arquivo)
	listenRPC("0.0.0.0:1313")
}