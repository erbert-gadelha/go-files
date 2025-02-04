package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

// func AbrirConexao(addr string) *net.UDPConn {
func AbrirConexao(addr string) net.Conn {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Printf("Erro ao Criar Conexão", err)
	}
	return conn
}

func handleConnection(conn net.Conn /**net.UDPConn*/) {
	params := [7]string{"olá\n", "mundo\n", "!\n", "como\n", "vai\n", "você\n", "?\n"}

	buffer := make([]byte, 1024)

	for i := 0; i < executions; i++ {
		fmt.Fprintf(conn, params[i%len(params)]+"\n")
		_, err := bufio.NewReader(conn).Read(buffer)
		if err == nil {
			fmt.Print(string(buffer))
		} else {
			fmt.Printf("Some error %v\n", err)
		}
	}
}

var executions = 10

func main() {
	if len(os.Args) > 1 {
		executions, _ = strconv.Atoi(os.Args[1])
	}

	conn := AbrirConexao("localhost:4002")
	defer conn.Close()

	handleConnection(conn)
}
