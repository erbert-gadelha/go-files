package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello from Server!")

	r, err := net.ResolveTCPAddr("tcp", "localhost:1313")
	if err != nil {
		fmt.Println("Erro ao resolver endereço:", err)
		os.Exit(1)
	}

	ln, err := net.ListenTCP("tcp", r)
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		os.Exit(1)
	}
	fmt.Println("Servidor iniciado na porta 1313.")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Cliente conectado:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler:", err)
			break
		}

		rep := strings.ToUpper(strings.TrimSpace(req))
		fmt.Println(req + " -> " + rep)

		_, err = conn.Write([]byte(rep + "\n"))
		if err != nil {
			fmt.Println("Erro ao enviar resposta:", err)
			break
		}
	}
	fmt.Println("Cliente desconectado.")
}
