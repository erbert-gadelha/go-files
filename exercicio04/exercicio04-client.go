package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"server"
)

func main() {
	fmt.Println("Hello from Client!")
	server.main();


	r, err := net.ResolveTCPAddr("tcp", "localhost:1313")
	if err != nil {
		fmt.Println("Erro ao resolver endereço:", err)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, r)
	if err != nil {
		fmt.Println("Erro ao conectar:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Print("Digite uma mensagem: ")
	reader := bufio.NewReader(os.Stdin)
	req, _ := reader.ReadString('\n')

	_, err = conn.Write([]byte(req))
	if err != nil {
		fmt.Println("Erro ao enviar mensagem:", err)
		os.Exit(1)
	}

	rep, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao receber resposta:", err)
		os.Exit(1)
	}

	fmt.Println("Servidor respondeu:", rep)
}
