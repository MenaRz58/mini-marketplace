package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"mini-marketplace/users/internal/controller/user"
	"mini-marketplace/users/internal/pkg/model"
)

type Handler struct {
	ctrl *user.Controller
}

func NewHandler(c *user.Controller) *Handler { return &Handler{ctrl: c} }

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.usersCollection)
	mux.HandleFunc("/users/", h.usersItem)
}

func (h *Handler) usersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(h.ctrl.List())
	case http.MethodPost:
		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		if u.ID == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		if err := h.ctrl.Create(u); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) usersItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		u, err := h.ctrl.Get(id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
