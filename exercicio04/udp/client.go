package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// func AbrirConexao(addr string) *net.UDPConn {
func AbrirConexao(addr string) net.Conn {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Printf("Erro ao Criar Conexão", err)
	}
	return conn
}

func handleConnection(conn net.Conn) {
	params := [7]string{"olá\n", "mundo\n", "!\n", "como\n", "vai\n", "você\n", "?\n"}

	buffer := make([]byte, 1024)

	for i := 0; i < executions; i++ {
		req := append([]byte{byte(i)}, []byte(params[i%len(params)])...)

		start := time.Now()

		fmt.Fprintf(conn, string(req))
		_, err := bufio.NewReader(conn).Read(buffer)
		if err != nil {
			fmt.Printf("Erro na resposta do servidor", err)
		}

		delta := time.Since(start) / time.Nanosecond
		fmt.Println(strconv.FormatInt(delta.Nanoseconds(), 10))
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
