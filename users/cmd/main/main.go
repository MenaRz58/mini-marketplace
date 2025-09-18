package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"mini-marketplace/pkg/discovery/consul"
	"mini-marketplace/users/internal/controller/user"
	uh "mini-marketplace/users/internal/handler/http"
	"mini-marketplace/users/internal/repository/memory"
)

func main() {
	serviceName := "users"
	port := getEnv("PORT", "8082")
	addr := "localhost:" + port

	repo := memory.NewUserRepository()
	registry, err := consul.NewRegistry("consul:8500")
	if err != nil {
		log.Fatal("Failed to connect to Consul:", err)
	}
	ctx := context.Background()
	instanceID := fmt.Sprintf("%s-%s", serviceName, port)

	// Registrar en Consul
	if err := registry.Register(ctx, instanceID, serviceName, addr); err != nil {
		log.Fatal("Failed to register service in Consul:", err)
	}
	defer registry.Deregister(ctx, instanceID, serviceName)
	log.Println("Service registered in Consul:", serviceName, "ID:", instanceID, "Address:", addr)

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state:", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	ctrl := user.NewController(repo)
	h := uh.NewHandler(ctrl)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	fmt.Println("Users service listening on http://" + addr)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
