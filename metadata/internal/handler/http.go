package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"mini-marketplace/metadata/internal/controller"
	metadataModel "mini-marketplace/metadata/pkg/model"
)

type Handler struct {
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/metadata", h.metadataHandler)
}

func (h *Handler) metadataHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		// GET requiere ?id=...
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		m, err := h.ctrl.Get(ctx, id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(m)

	case http.MethodPost:
		// POST recibe JSON con los campos de Metadata
		var meta metadataModel.Metadata
		if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		// Generar ID si no lo mandan
		if meta.ID == "" {
			meta.ID = uuid.New().String()
		}
		if err := h.ctrl.Put(ctx, meta.ID, &meta); err != nil {
			http.Error(w, "cannot save", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
