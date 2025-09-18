package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"mini-marketplace/pkg/discovery/consul"
	"mini-marketplace/products/internal/controller/product"
	phttp "mini-marketplace/products/internal/handler/http"
	"mini-marketplace/products/internal/repository/memory"
)

func main() {
	serviceName := "products"
	port := getEnv("PORT", "8081")
	addr := "localhost:" + port

	// Repositorio y controlador
	repo := memory.NewProductRepository()
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

	ctrl := product.NewController(repo)
	h := phttp.NewHandler(ctrl)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	fmt.Println("Products service listening on http://" + addr)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
