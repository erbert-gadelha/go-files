package main

import "sync"

type Semaforo struct {
    n int32
    ch chan int32
}

type Banco struct {
    saldo int32
	semaforo *Semaforo
}

func criarSemaforo(N int32) *Semaforo {
    return &Semaforo{
        N,
        make(chan int32, N),
    }
}


func criarBanco(saldo int32, n int32) *Banco {
    return &Banco{
        saldo,
        criarSemaforo(n),
    }
}

func (s *Semaforo) P() {
	s.ch<-0
}

func (s *Semaforo) V() {
	<-s.ch
}

func (b *Banco) depositar(valor int32) {
	b.semaforo.P()
	b.saldo+=valor;
	b.semaforo.V()
}

func (b *Banco) sacar(valor int32) {
	b.semaforo.P()
	b.saldo-=valor;
	b.semaforo.V()
}

func main() {

	var banco *Banco = criarBanco(0, 1)
    var wg sync.WaitGroup
    wg.Add(3)

    go func(banco *Banco) {
        defer wg.Done()
        for i := 0; i < 8096; i++ {
			banco.depositar(1);
        }
    }(banco)

    go func(banco *Banco) {
        defer wg.Done()
        for i := 0; i < 8096; i++ {
			banco.sacar(1);
        }
    }(banco)


    go func(banco *Banco) {
        defer wg.Done()
        for i := 0; i < 8096; i++ {
			banco.depositar(1);
        }
    }(banco)

    wg.Wait()
    print("saldo: ", banco.saldo, ".\n")
}