package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Higor-ViniciusDev/auth/configuration/rest_err"
	user_usecase "github.com/Higor-ViniciusDev/auth/internal/usecase/user"
)

type AuthLoginHandler struct {
	userUseCase *user_usecase.UserUseCase
}

func NewAuthLoginHandler(userUseCase *user_usecase.UserUseCase) *AuthLoginHandler {
	return &AuthLoginHandler{userUseCase: userUseCase}
}

func (h *AuthLoginHandler) Autenticacao(w http.ResponseWriter, r *http.Request) {
	var dto user_usecase.UserInputValidationDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	errInternal := h.userUseCase.Validation(ctx, dto)

	if errInternal != nil {
		restErro := rest_err.ConvertInternalErrorToRestError(errInternal)
		w.WriteHeader(restErro.Code)
		json.NewEncoder(w).Encode(restErro)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]byte("ola mundo"))
}

func (h *AuthLoginHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto user_usecase.CreateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	errInternal := h.userUseCase.Create(ctx, dto)

	if errInternal != nil {
		restErro := rest_err.ConvertInternalErrorToRestError(errInternal)
		w.WriteHeader(restErro.Code)
		json.NewEncoder(w).Encode(restErro)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
