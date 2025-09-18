package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"mini-marketplace/pkg/discovery/registry"
)

type Handler struct {
	reg registry.Discovery
}

func New(reg registry.Discovery) *Handler {
	return &Handler{reg: reg}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users/", h.proxy("users"))
	mux.HandleFunc("/products/", h.proxy("products"))
	mux.HandleFunc("/orders/", h.proxy("orders"))
	mux.HandleFunc("/metadata/", h.proxy("metadata"))
}

func (h *Handler) proxy(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instances, err := h.reg.ServiceAddress(r.Context(), serviceName)
		if err != nil || len(instances) == 0 {
			http.Error(w, "service not available: "+serviceName, http.StatusServiceUnavailable)
			return
		}

		// Construir URL del servicio
		target := &url.URL{
			Scheme: "http",
			Host:   instances[0],
		}

		proxy := httputil.NewSingleHostReverseProxy(target)

		// Ajustar Host para el backend
		r.Host = instances[0]

		log.Println("Proxying request to:", target.String()+r.URL.Path)
		proxy.ServeHTTP(w, r)
	}
}
