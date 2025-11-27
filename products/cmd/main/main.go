package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"google.golang.org/grpc"

	"mini-marketplace/products/internal/controller/product"
	"mini-marketplace/products/internal/pkg/model"
	productRepo "mini-marketplace/products/internal/repository/postgres"
	"mini-marketplace/products/internal/service"
	pb "mini-marketplace/proto/products"
)

func main() {
	port := "50051"
	addr := fmt.Sprintf(":%s", port)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fallo DB:", err)
	}

	// 2. Migración
	db.AutoMigrate(&model.Product{})

	// 3. Repo SQL
	repo := productRepo.NewProductRepository(db)
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
