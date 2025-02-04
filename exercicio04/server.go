package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

func AbrirConexaoTCP(addr string) *net.TCPListener {
	fmt.Println("Hello from Server!")

	r, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("Erro ao resolver endereço:", err)
		os.Exit(1)
	}

	ln, err := net.ListenTCP("tcp", r)
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		os.Exit(1)
	}

	return ln
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

func CriarArquivos() {
	for i := 0; i < len(arquivos); i++ {
		arquivos[i] = CriarArquivo(fmt.Sprintf("exercicio04-ebgr-%02d.txt", i))
	}
}

var arquivos [4]*Arquivo

func main() {
	ln := AbrirConexaoTCP("localhost:1313")
	fmt.Println("Servidor iniciado na porta 1313.")

	CriarArquivos()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		go handleConnection_(conn)
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

// ///// SUBROTINA QUE REALIZA CONCATENACAO (LEITURA E ESCRITA) NO ARQUIVO	/////////
func EscreverArquivos_WAITGROUP(arquivos []*Arquivo, nome string, count int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		for j := 0; j < len(arquivos); j++ {
			arquivos[j].Concatenar(fmt.Sprintf("%s: %d\n", nome, i))
		}
	}
}

// ///// SUBROTINAS CONCORRENTES	/////////
func func_parallel(arquivos []*Arquivo, count int) {
	var wg sync.WaitGroup
	wg.Add(5)

	go EscreverArquivos_WAITGROUP(arquivos[:], "thread A", count, &wg)
	go EscreverArquivos_WAITGROUP(arquivos[:], "thread B", count, &wg)
	go EscreverArquivos_WAITGROUP(arquivos[:], "thread C", count, &wg)
	go EscreverArquivos_WAITGROUP(arquivos[:], "thread D", count, &wg)
	go EscreverArquivos_WAITGROUP(arquivos[:], "thread E", count, &wg)
	wg.Wait()
}
