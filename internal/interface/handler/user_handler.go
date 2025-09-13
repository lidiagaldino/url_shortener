package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/services"
	"url-shortener/internal/services/dto"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Save(w http.ResponseWriter, r *http.Request) {
	var req dto.UserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	user, err := h.service.Save(&req)
	if err != nil {
		http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	user, err := h.service.LoginUser(&req)
	if err != nil {
		if errors.Is(err, exceptions.ErrInvalidCredentials) {
			http.Error(w, "Credenciais invalidas", http.StatusForbidden)
			return
		}
		http.Error(w, "Erro ao fazer login", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
