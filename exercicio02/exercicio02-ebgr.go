///////////////////////////////////////////////////////
// Equipe 09 - ERBERT BERNARDINO GADELHA ROCHA[ebgr] //
//                                                   //
//      Implementar um programa que lê e escreve     //
//     em múltiplos arquivos concorrentemente.       //
//                                                   //
///////////////////////////////////////////////////////
//	NA PRIMEIRA ENTREGA A APLICAÇÃO NÃO LIDAVA COM   //
//	MULTIPLOS ARQUIVOS. ISSO FOI CORRIGIDO AQUI.     //
//	                                                 //
//	O MONITOR E SEMÁFORO SEGUEM INALTERADOS.         //
///////////////////////////////////////////////////////

package main

import (
	"fmt"
	"os"
	"sync"
    "time"
	"math"
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
	func EscreverArquivos_WAITGROUP(arquivos []*Arquivo, nome string, count int, wg *sync.WaitGroup) {
		defer wg.Done()
		for i:=0; i < count; i++ {
			for j := 0; j < len(arquivos); j++ {
				arquivos[j].Concatenar(fmt.Sprintf("%s: %d\n", nome, i))
			}
		}
	}

	func EscreverArquivos_NO_WAITGROUP(arquivos []*Arquivo, nome string, count int) {
		for i:=0; i < count; i++ {
			for j := 0; j < len(arquivos); j++ {
				arquivos[j].Concatenar(fmt.Sprintf("%s: %d\n", nome, i))
			}
		}
	}




/////// SUBROTINAS CONCORRENTES	/////////
	func func_parallel() {
		var arquivos [10]*Arquivo;
		for i:=0; i<len(arquivos); i++ {
			arquivos[i] = CriarArquivo(fmt.Sprintf("exercicio02-ebgr-%02d.txt", i))
		}

		var wg sync.WaitGroup
		wg.Add(5)
		var count int = 10

		go EscreverArquivos_WAITGROUP(arquivos[:], "thread A", count, &wg)
		go EscreverArquivos_WAITGROUP(arquivos[:], "thread B", count, &wg)
		go EscreverArquivos_WAITGROUP(arquivos[:], "thread C", count, &wg)
		go EscreverArquivos_WAITGROUP(arquivos[:], "thread D", count, &wg)
		go EscreverArquivos_WAITGROUP(arquivos[:], "thread E", count, &wg)	
		wg.Wait()
	}

/////// SUBROTINAS SEQUENCIAIS	/////////
	func func_synchronous() {
		var arquivos [10]*Arquivo;
		for i:=0; i<len(arquivos); i++ {
			arquivos[i] = CriarArquivo(fmt.Sprintf("exercicio02-ebgr-%02d.txt", i))
		}

		var count int = 10
		EscreverArquivos_NO_WAITGROUP(arquivos[:], "thread A", count)
		EscreverArquivos_NO_WAITGROUP(arquivos[:], "thread B", count)
		EscreverArquivos_NO_WAITGROUP(arquivos[:], "thread C", count)
		EscreverArquivos_NO_WAITGROUP(arquivos[:], "thread D", count)
		EscreverArquivos_NO_WAITGROUP(arquivos[:], "thread E", count)
	}

/////// MEDE TEMPO DE EXECUCAO	/////////
	func measure(name string, method func(), attempts int) float64 {
		var mediam float64 = 0;
		var deviation float64 = 0;

		fmt.Printf("[%s]:\n", name)
		for i:=0; i<attempts; i++ {
			start := time.Now()
			method()
			delta := float64(time.Since(start)) / float64(time.Millisecond)
			fmt.Printf("[%02d/%d]: %f\n", i, attempts, delta)
			mediam += delta
			deviation += math.Pow(delta, 2)
		}

		mediam = mediam/float64(attempts)
		deviation = math.Sqrt(deviation/float64(attempts))
		fmt.Printf("%s: [média=%fms, desvio=%fms]\n\n", name, mediam, deviation)
		return mediam;
	}

func main() {
	var attempts int = 10;
	measure("SYNCHRONOUS", func_synchronous, attempts)
	measure(   "PARALLEL",    func_parallel, attempts)
}


