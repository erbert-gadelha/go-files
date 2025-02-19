package main

import (
	"context"
	"log"
	"time"

	pb "exercicio05/grpc/arquivo" // Importe o pacote gerado

	"google.golang.org/grpc"
)

func main() {
	// Conectar ao servidor gRPC
	conn, err := grpc.Dial("localhost:1313", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Erro ao conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewArquivoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Enviar uma mensagem para o servidor
	req := &pb.Args{Conteudo: "Exemplo de conteúdo"}
	resp, err := client.Linhas(ctx, req)
	if err != nil {
		log.Fatalf("Erro ao chamar Linhas: %v", err)
	}

	log.Printf("Resposta do servidor: %d", resp.Linhas)
}
