package main

import (
	"context"
	pb "exercicio05/grpc/arquivo" // Certifique-se de que este caminho está correto
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedArquivoServiceServer
}

func (s *server) Linhas(ctx context.Context, req *pb.Args) (*pb.Reply, error) {
	log.Printf("Recebido: %s", req.Conteudo)
	return &pb.Reply{Linhas: int32(len(req.Conteudo))}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":1313")
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterArquivoServiceServer(grpcServer, &server{})

	log.Println("Servidor gRPC rodando na porta 1313...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Falha ao iniciar servidor: %v", err)
	}
}
