package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mini-marketplace/orders/internal/controller/order"
	"mini-marketplace/orders/internal/pkg/model"
	orderRepo "mini-marketplace/orders/internal/repository/postgres"
	"mini-marketplace/orders/internal/service"
	pb "mini-marketplace/proto/orders"
)

func main() {
	port := 50052
	addr := fmt.Sprintf(":%d", port)

	// Repositorio en memoria
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fallo al conectar a la DB:", err)
	}

	err = db.AutoMigrate(&model.Order{}, &model.OrderProduct{})
	if err != nil {
		log.Fatal("Fallo en la migración:", err)
	}

	repo := orderRepo.NewOrderRepository(db)

	ctrl := order.NewController(repo)
	grpcHandler := service.New(ctrl)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrdersServiceServer(grpcServer, grpcHandler)

	log.Printf("Orders gRPC server listening on %s", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
