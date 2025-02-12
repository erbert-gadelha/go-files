package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func OuvirConexaoUDP(port int) *net.UDPConn {
	config := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
		Zone: "",
	}

	conn, err := net.ListenUDP("udp", config)
	if err != nil {
		fmt.Println("Erro ao resolver endereço:", err)
		os.Exit(1)
	}

	fmt.Printf("Hello from UDP Server!\nServidor iniciado na porta %d.\n", port)

	return conn
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

var arquivos []*Arquivo

func handleMessageUDP(conn *net.UDPConn, remoteAddr net.Addr, buf []byte) {

	index := int(buf[0])
	buf = buf[1:]
	newlineIndex := bytes.IndexByte(buf, '\n')
	if newlineIndex != -1 {
		buf = buf[:newlineIndex]
	}

	var rep string = strconv.Itoa(EscreverArquivo(string(buf)+"\n", index))
	_, err := conn.WriteTo([]byte(rep+"\n"), remoteAddr)
	if err != nil {
		fmt.Println("Erro ao enviar resposta", err)
	}
}

func main() {
	CriarArquivos(10)

	conn := OuvirConexaoUDP(4002)
	defer conn.Close()

	for {
		buf := make([]byte, 2048)
		_, remoteAddr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Cliente conectado:", remoteAddr)

		go handleMessageUDP(conn, remoteAddr, buf)
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
