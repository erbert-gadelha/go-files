package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func AbrirConexao(addr string) *net.TCPConn {
	r, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("Erro ao resolver endereço:", err)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, r)
	if err != nil {
		fmt.Println("Erro ao conectar:", err)
		os.Exit(1)
	}

	return conn
}

func handleConnection(conn *net.TCPConn) {
	params := [7]string{"olá\n", "mundo\n", "!\n", "como\n", "vai\n", "você\n", "?\n"}

	for i := 0; i < executions; i++ {
		req := append([]byte{byte(i)}, []byte(params[i%len(params)])...)

		start := time.Now()
		_, err := conn.Write(req)
		if err != nil {
			fmt.Println("Erro ao enviar mensagem:", err)
			os.Exit(1)
		}

		_, err = bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao receber resposta:", err)
			os.Exit(1)
		}

		//rep = rep[:len(rep)-1]
		delta := time.Since(start) / time.Nanosecond
		fmt.Println(strconv.FormatInt(delta.Nanoseconds(), 10))
	}
}

var executions = 10

func main() {
	if len(os.Args) > 1 {
		executions, _ = strconv.Atoi(os.Args[1])
	}

	conn := AbrirConexao("localhost:1313")
	defer conn.Close()

	handleConnection(conn)
}
