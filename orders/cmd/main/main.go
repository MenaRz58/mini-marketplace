package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"mini-marketplace/orders/internal/controller/order"
	oh "mini-marketplace/orders/internal/handler/http"
	"mini-marketplace/orders/internal/repository/memory"
	"mini-marketplace/pkg/discovery/consul"
)

func main() {
	serviceName := "orders"
	port := getEnv("PORT", "8083")
	addr := "localhost:" + port

	// Repositorio
	repo := memory.NewInMemoryRepository()

	// Conectamos a Consul
	registry, err := consul.NewRegistry("consul:8500")
	if err != nil {
		log.Fatal("Failed to connect to Consul:", err)
	}

	ctx := context.Background()
	instanceID := fmt.Sprintf("%s-%s", serviceName, port)

	// Registramos el servicio
	if err := registry.Register(ctx, instanceID, serviceName, addr); err != nil {
		log.Fatal("Failed to register service in Consul:", err)
	}
	defer registry.Deregister(ctx, instanceID, serviceName)
	log.Println("Service registered in Consul:", serviceName, "ID:", instanceID, "Address:", addr)

	// Goroutine para reportar estado saludable
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state:", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Controlador
	ctrl := order.NewController(repo, registry, ctx)
	h := oh.NewHandler(ctrl)

	// Router y rutas
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	fmt.Println("Orders service listening on http://" + addr)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
