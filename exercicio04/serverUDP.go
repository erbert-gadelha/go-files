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

func AbrirConexaoUDP(port int) *net.UDPConn {
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

	return conn
}

func handleConnection_(conn net.Conn) {
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

var arquivos []*Arquivo

func handleMessage(conn *net.UDPConn, remoteAddr net.Addr, buf []byte) {

	str := string(buf)
	for i := 0; i < len(buf); i++ {
		if buf[i] == '\n' {
			str = str[:i]
			break
		}
	}
	rep := strconv.Itoa(len(str))
	_, err := conn.WriteTo([]byte(rep+"\n"), remoteAddr)
	if err != nil {
		fmt.Println("Erro ao enviar resposta", err)
	}

	fmt.Println(fmt.Sprintf("[%s]: \"%s\"", remoteAddr, str))
}

func main() {
	CriarArquivos(2)

	conn := AbrirConexaoUDP(4002)
	defer conn.Close()

	//remoteConns := new(sync.Map)

	for {
		buf := make([]byte, 1024)
		_, remoteAddr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		handleMessage(conn, remoteAddr, buf)
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
