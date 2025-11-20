package main

import (
	"fmt"
	"log"
	"net"

	"mini-marketplace/orders/internal/controller/order"
	"mini-marketplace/orders/internal/repository/memory"
	pb "mini-marketplace/proto/orders"

	"google.golang.org/grpc"
)

func main() {
	port := 50052
	addr := fmt.Sprintf(":%d", port)

	// Repositorio en memoria
	repo := memory.NewInMemoryRepository()

	// Controlador
	ctrl := order.NewController(repo)

	// gRPC Server
	grpcServer := grpc.NewServer()
	pb.RegisterOrdersServiceServer(grpcServer, order.NewGRPCServer(ctrl))

	log.Printf("Orders gRPC server listening on %s", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
