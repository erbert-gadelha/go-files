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
	"io"
	"math"
	"os"
	"sync"
	"time"
)

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
	escrita  *os.File
	leitura  *os.File
	semaforo Semaforo
}

func CriarArquivo(nome string) (a *Arquivo) {
	leitura, _ := os.Open(nome)
	content, _ := io.ReadAll(leitura)

	escrita, _ := os.Create(nome)
	escrita.WriteString(string(content))

	return &Arquivo{nome, escrita, leitura, Semaforo{make(chan bool, 1)}}
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

func EscreverArquivos_NO_WAITGROUP(arquivos []*Arquivo, nome string, count int) {
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

// ///// SUBROTINAS SEQUENCIAIS	/////////
func func_synchronous(arquivos []*Arquivo, count int) {
	EscreverArquivos_NO_WAITGROUP(arquivos[:], "thread A", count)
	EscreverArquivos_NO_WAITGROUP(arquivos[:], "thread B", count)
	EscreverArquivos_NO_WAITGROUP(arquivos[:], "thread C", count)
	EscreverArquivos_NO_WAITGROUP(arquivos[:], "thread D", count)
	EscreverArquivos_NO_WAITGROUP(arquivos[:], "thread E", count)
}

// ///// MEDE TEMPO DE EXECUCAO	/////////
func measure(name string, method func(arquivos []*Arquivo, count int), attempts int) float64 {
	var mediam float64 = 0
	var deviation float64 = 0

	var arquivos [10]*Arquivo
	for i := 0; i < len(arquivos); i++ {
		arquivos[i] = CriarArquivo(fmt.Sprintf("exercicio02-ebgr-%02d.txt", i))
	}

	var times []float64 = make([]float64, attempts)

	fmt.Printf("[%s]:\n", name)
	attempts--
	for i := 0; i <= attempts; i++ {
		start := time.Now()
		method(arquivos[:], 100)
		delta := float64(time.Since(start)) / float64(time.Millisecond)
		times[i] = delta
		fmt.Printf("[%03d/%03d]: %f\n", i, attempts, delta)
		mediam += delta
	}

	mediam /= float64(attempts)
	for i := 0; i < attempts; i++ {
		difference := mediam - times[i]
		deviation += (difference * difference)
	}

	deviation = math.Sqrt(deviation / float64(attempts))

	fmt.Printf("%s: [média=%fms, desvio=%fms]\n\n", name, mediam, deviation)
	return mediam
}

func main() {
	var attempts int = 20
	measure("SYNCHRONOUS", func_synchronous, attempts)
	measure("PARALLEL", func_parallel, attempts)
}
