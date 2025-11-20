package main

import (
	"fmt"
	"log"
	"net"

	pb "mini-marketplace/proto/metadata"

	"mini-marketplace/metadata/internal/controller"
	"mini-marketplace/metadata/internal/repository"
	"mini-marketplace/metadata/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 50051

	// Dependencias
	repo := repository.New()
	ctrl := controller.New(repo)
	server := service.New(ctrl)

	// Servidor gRPC
	grpcServer := grpc.NewServer()

	// Registrar servicio
	pb.RegisterMetadataServiceServer(grpcServer, server)

	// Reflection (para grpcurl)
	reflection.Register(grpcServer)

	// Escuchar puerto
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Metadata gRPC server running on port %d...", port)

	// Iniciar
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
