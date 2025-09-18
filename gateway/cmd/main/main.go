package main

import (
	"fmt"
	"log"
	"net/http"

	"mini-marketplace/gateway/internal/handler"
	"mini-marketplace/pkg/discovery/consul"
)

func main() {
	registry, err := consul.NewRegistry("consul:8500")
	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(registry)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	port := 8080
	fmt.Println("Gateway listening on", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
