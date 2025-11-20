package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"mini-marketplace/products/internal/controller/product"
	"mini-marketplace/products/internal/repository/memory"
	"mini-marketplace/products/internal/service"
	pb "mini-marketplace/proto/products"
)

func main() {
	port := "50051"
	addr := fmt.Sprintf(":%s", port)

	// Repositorio y controlador
	repo := memory.NewProductRepository()
	ctrl := product.NewController(repo)

	// Servidor gRPC
	grpcServer := service.New(ctrl)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	s := grpc.NewServer()
	pb.RegisterProductsServiceServer(s, grpcServer)

	log.Println("Products gRPC server listening on", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
