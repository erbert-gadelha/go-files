package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var arquivos []*Arquivo

func AbrirConexaoTCP(port int) *net.TCPListener {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Println("Erro ao resolver endereço:", err)
		os.Exit(1)
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		os.Exit(1)
	}

	fmt.Printf("Hello from TCP Server!\nServidor iniciado na porta (%d).", port)

	return ln
}

func handleConnectionTCP(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Cliente conectado:", conn.RemoteAddr())
	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		var rep string = strconv.Itoa(EscreverArquivo(req[1:], int(req[0])))
		_, err = conn.Write([]byte(rep + "\n"))
		if err != nil {
			break
		}
	}

	fmt.Println("Cliente desconectado.", conn.RemoteAddr())
}

func EscreverArquivo(str string, index int) int {
	if index < 0 {
		return -1
	}
	arquivos[index%len(arquivos)].Concatenar(str)
	return arquivos[index%len(arquivos)].linhas
}

func CriarArquivos(count int) {
	arquivos = make([]*Arquivo, count)

	for i := 0; i < len(arquivos); i++ {
		arquivos[i] = CriarArquivo(fmt.Sprintf("exercicio04-ebgr-%02d.txt", i))
	}
}

func main() {
	ln := AbrirConexaoTCP(1313)

	CriarArquivos(10)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleConnectionTCP(conn)
	}
}

// ///// SEMAFORO IMPLEMENTADO COM CANAL	/////////
type Semaforo struct {
	ch chan bool
}

func (s *Semaforo) P() {
	s.ch <- false
}
func (s *Semaforo) V() {
	<-s.ch
}

// ///// MONITOR QUE GARANTE ATOMICIDADE DAS OPERACOES	/////////
type Arquivo struct {
	nome     string
	linhas   int
	escrita  *os.File
	leitura  *os.File
	semaforo Semaforo
}

func CriarArquivo(nome string) (a *Arquivo) {
	leitura, _ := os.Open(nome)
	content, _ := io.ReadAll(leitura)

	str := string(content)

	escrita, _ := os.Create(nome)
	escrita.WriteString(str)

	linhas := strings.Count(str, "\n")

	return &Arquivo{nome, linhas, escrita, leitura, Semaforo{make(chan bool, 1)}}
}

func (a *Arquivo) Ler() (result string) {
	a.semaforo.P()
	defer a.semaforo.V()
	content, _ := io.ReadAll(a.leitura)
	return string(content)
}

func (a *Arquivo) Concatenar(value string) {
	a.semaforo.P()
	defer a.semaforo.V()
	a.escrita.WriteString(value)
	a.linhas++
}
