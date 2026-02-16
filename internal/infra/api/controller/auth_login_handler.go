package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	w.Write([]byte(fmt.Sprintf("Login: %v | Senha: %v", dto.Email, dto.Password)))

}
