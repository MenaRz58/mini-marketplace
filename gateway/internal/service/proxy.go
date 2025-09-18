package service

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"mini-marketplace/pkg/discovery/consul"
)

func NewProxy(registry *consul.Registry, serviceName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Buscar instancias en Consul
		instances, err := registry.ServiceAddress(r.Context(), serviceName)
		if err != nil || len(instances) == 0 {
			http.Error(w, "service not available", http.StatusServiceUnavailable)
			return
		}

		// Usar la primera instancia
		targetURL, _ := url.Parse("http://" + instances[0])

		// Proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(w, r)
	})
}
