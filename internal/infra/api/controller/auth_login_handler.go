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
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	userDto, err := h.userUseCase.Validation(ctx, dto)
	if err != nil {
		restErro := rest_err.ConvertInternalErrorToRestError(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(restErro.Code)
		json.NewEncoder(w).Encode(restErro)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userDto)
}

func (h *AuthLoginHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto user_usecase.CreateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	errInternal := h.userUseCase.Create(ctx, dto)
	if errInternal != nil {
		restErro := rest_err.ConvertInternalErrorToRestError(errInternal)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(restErro.Code)
		json.NewEncoder(w).Encode(restErro)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthLoginHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	errInternal := h.userUseCase.VerifyEmail(ctx, token)
	if errInternal != nil {
		restErro := rest_err.ConvertInternalErrorToRestError(errInternal)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(restErro.Code)
		json.NewEncoder(w).Encode(restErro)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "email verified successfully"})
}

func (h *AuthLoginHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var dto user_usecase.ResendVerificationDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	errInternal := h.userUseCase.ResendVerification(ctx, dto)
	if errInternal != nil {
		restErro := rest_err.ConvertInternalErrorToRestError(errInternal)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(restErro.Code)
		json.NewEncoder(w).Encode(restErro)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "email resent"})
}
