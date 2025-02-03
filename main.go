package main

import "sync"

type Semaforo struct {
    n int32
    c *sync.Cond
}

type Banco struct {
    saldo int32
}

func criarSemaforo(N int32) *Semaforo {
    return &Semaforo{
        N,
        sync.NewCond(new(sync.Mutex)),
    }
}

func (s *Semaforo) P() {
    s.c.L.Lock()
    for s.n <= 0 {
        s.c.Wait()
    }
    s.n--
    s.c.L.Unlock()
}

func (s *Semaforo) V() {
    s.c.L.Lock()
    s.n++
    s.c.L.Unlock()
    s.c.Signal()
}

func (b *Banco) depositar(valor int32, s Semaforo) {
    s.P()
    b.saldo += valor
    s.V()
}

func (b *Banco) sacar(valor int32, s Semaforo) {
    s.P()
    b.saldo -= valor
    s.V()
}

func main() {

    var ch = make(chan int, 1)
    var saldo int32 = 0
    var wg sync.WaitGroup
    wg.Add(2)

    go func(ch chan int) {
        defer wg.Done()
        for i := 0; i < 8096; i++ {
            ch <- 1
            saldo += 1
            <-ch
        }
    }(ch)

    go func(ch chan int) {
        defer wg.Done()
        for i := 0; i < 8096; i++ {
            ch <- 1
            saldo -= 1
            <-ch
        }
    }(ch)

    wg.Wait()
    print("saldo: ", saldo, ".\n")
}