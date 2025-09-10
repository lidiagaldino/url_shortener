package handler

import (
	"encoding/json"
	"net/http"
	"url-shortener/internal/services"
	"url-shortener/internal/services/dto"

	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	service *services.URLService
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{service: service}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req dto.URL
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	url, err := h.service.Shorten(req.URL)
	if err != nil {
		http.Error(w, "Erro ao encurtar", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(url)
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, err := h.service.Resolve(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
}
