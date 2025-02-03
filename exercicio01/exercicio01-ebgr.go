//////////////////////////////////////////////////////
//   Equipe 09 | ERBERT BERNARDINO GADELHA ROCHA    //
//                                                  //
//      Implementar um programa que lê e escreve    //
//     em múltiplos arquivos concorrentemente.      //
//                                                  //
//////////////////////////////////////////////////////
//	CONCORRENCIA SOBRE O ARQUIVO RESOLVIDA COM		//
//	CANAL DE TAMANHO ÚNICO.							//
//	GOROTINA PRINCIPAL AGUARDA SUBROTINAS POR		//
//	MEIO DE WAITGROUP.								//
//////////////////////////////////////////////////////

package main

import (
	"fmt"
	"os"
	"sync"
)

/////// SEMAFORO IMPLEMENTADO COM CANAL	/////////
	type Semaforo struct {
		ch chan bool
	}

	func (s *Semaforo) P() {
		s.ch <- false
	}
	func (s *Semaforo) V() {
		<-s.ch
	}

/////// MONITOR QUE GARANTE ATOMICIDADE DAS OPERACOES	/////////
	type Arquivo struct {
		nome     string
		semaforo Semaforo
	}

	func AbrirArquivo(nome string) (a *Arquivo) {
		return &Arquivo{nome, Semaforo{make(chan bool, 1)}}
	}

	func CriarArquivo(nome string) (a *Arquivo) {
		file, err := os.Create(nome)
		defer file.Close()
		if err != nil {
			panic(err)
		}
		return &Arquivo{nome, Semaforo{make(chan bool, 1)}}
	}

	func (a *Arquivo) Ler() (result string) {
		a.semaforo.P()
		defer a.semaforo.V()

		content, err := os.ReadFile(a.nome)
		if err != nil {
			panic(err)
		}

		return string(content)
	}

	func (a *Arquivo) Escrever(value string) {
		a.semaforo.P()
		defer a.semaforo.V()

		file, err := os.Create(a.nome)
		defer file.Close()

		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(value)
		if err != nil {
			panic(err)
		}
	}

	func (a *Arquivo) Concatenar(value string) {
		a.semaforo.P()
		defer a.semaforo.V()

		content, err := os.ReadFile(a.nome)
		if err != nil {
			panic(err)
		}

		file, err := os.Create(a.nome)
		defer file.Close()

		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(string(content) + value)
		if err != nil {
			panic(err)
		}
	}

/////// SUBROTINA QUE REALIZA CONCATENACAO (LEITURA E ESCRITA) NO ARQUIVO	/////////
	func EscreverArquivo(arquivo *Arquivo, nome string, count int, wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < count; i++ {
			arquivo.Concatenar(fmt.Sprintf("%s: %d\n", nome, i))
		}
	}

func main() {
	var arquivo *Arquivo = AbrirArquivo("saida.txt")
	arquivo.Escrever("")

	var wg sync.WaitGroup
	wg.Add(5)

	var count int = 100
	go EscreverArquivo(arquivo, "thread A", count, &wg)
	go EscreverArquivo(arquivo, "thread B", count, &wg)
	go EscreverArquivo(arquivo, "thread C", count, &wg)
	go EscreverArquivo(arquivo, "thread D", count, &wg)
	go EscreverArquivo(arquivo, "thread E", count, &wg)

	wg.Wait()
}

