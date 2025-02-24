package main

import (
	"context"
	"log"
	"net"

	//"time"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	pb "server.go/grpcarquivo"
)

type server struct {
	pb.UnimplementedArquivoServiceServer
}

func (s *server) CountLines(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Lines: 1 + int32(strings.Count(req.GetContent(), "\n"))}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":1313")
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterArquivoServiceServer(s, &server{})

	fmt.Println("Servidor gRPC rodando na porta 1313...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Erro ao servir: %v", err)
	}
}
