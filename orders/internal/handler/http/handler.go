package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"mini-marketplace/orders/internal/controller/order"
	"mini-marketplace/orders/internal/pkg/model"
)

type Handler struct {
	ctrl *order.Controller
}

func NewHandler(c *order.Controller) *Handler { return &Handler{ctrl: c} }

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/orders", h.ordersCollection)
	mux.HandleFunc("/orders/", h.ordersItem)
}

func (h *Handler) ordersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(h.ctrl.List())
	case http.MethodPost:
		var o model.Order
		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		if err := h.ctrl.Create(o); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) ordersItem(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/orders/")
	if path == "" {
		// Si no hay ID, devuelve lista
		h.ordersCollection(w, r)
		return
	}
	id := path
	switch r.Method {
	case http.MethodGet:
		o, err := h.ctrl.Get(id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(o)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
