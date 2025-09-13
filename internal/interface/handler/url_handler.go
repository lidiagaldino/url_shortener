package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/interface/middleware"
	"url-shortener/internal/services"
	"url-shortener/internal/services/dto"

	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	service *services.URLService
}

func NewURLHandler(s *services.URLService) *URLHandler {
	return &URLHandler{service: s}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.URL
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	url, err := h.service.Shorten(req.URL, userID)
	if err != nil {
		if errors.Is(err, exceptions.ErrInvalidURL) {
			http.Error(w, "URL inválida", http.StatusBadRequest)
			return
		}
		http.Error(w, "Erro ao encurtar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(url)
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ip := r.RemoteAddr
	userAgent := r.UserAgent()
	referer := r.Referer()

	url, err := h.service.Resolve(id, ip, userAgent, referer)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
}
