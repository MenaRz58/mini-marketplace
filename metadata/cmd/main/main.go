package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"mini-marketplace/metadata/internal/controller"
	"mini-marketplace/metadata/internal/handler"
	"mini-marketplace/metadata/internal/repository"
	"mini-marketplace/pkg/discovery/consul"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8090, "metadata API port")
	flag.Parse()

	registry, err := consul.NewRegistry("consul:8500")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	instanceID := fmt.Sprintf("metadata-%d", port)
	if err := registry.Register(ctx, instanceID, "metadata", fmt.Sprintf("localhost:%d", port)); err != nil {
		log.Fatal(err)
	}
	defer registry.Deregister(ctx, instanceID, "metadata")

	// health
	go func() {
		for {
			registry.ReportHealthyState(instanceID, "metadata")
			time.Sleep(5 * time.Second)
		}
	}()

	repo := repository.New()
	ctrl := controller.New(repo) // <- ahora usamos el nombre real del paquete
	h := handler.New(ctrl)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	fmt.Println("Metadata service listening on", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
