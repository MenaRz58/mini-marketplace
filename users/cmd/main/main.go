package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "mini-marketplace/proto/users"
	adminpb "mini-marketplace/proto/users/admin"

	"mini-marketplace/users/internal/controller/user"
	"mini-marketplace/users/internal/repository/memory"
	"mini-marketplace/users/internal/service"
)

func main() {
	// 1. Crear repo en memoria
	repo := memory.NewUserRepository()

	// 2. Crear controlador (Ya no necesita clientes externos)
	ctrl, err := user.NewController(repo)
	if err != nil {
		log.Fatalf("failed to create controller: %v", err)
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
