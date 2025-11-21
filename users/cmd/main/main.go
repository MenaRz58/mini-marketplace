package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"google.golang.org/grpc"

	pb "mini-marketplace/proto/users"
	adminpb "mini-marketplace/proto/users/admin"

	"mini-marketplace/users/internal/controller/user"
	"mini-marketplace/users/internal/pkg/model"
	userRepo "mini-marketplace/users/internal/repository/postgres"
	"mini-marketplace/users/internal/service"
)

func main() {
	// 1. Crear repo en memoria
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

	db.AutoMigrate(&model.User{})

	repo := userRepo.NewUserRepository(db)

	// 2. Crear controlador
	ctrl, err := user.NewController(repo)
	if err != nil {
		log.Fatalf("Error creando el controlador: %v", err)
	}

	// 3. Crear servicios gRPC
	usersSvc := service.New(ctrl)
	adminSvc := service.NewAdmin(ctrl)

	// 4. Levantar servidor
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Registrar ambos servicios
	pb.RegisterUsersServiceServer(grpcServer, usersSvc)
	adminpb.RegisterUsersAdminServiceServer(grpcServer, adminSvc)

	log.Println("Users gRPC server listening on :50054")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
