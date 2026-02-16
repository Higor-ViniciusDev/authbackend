package user_usecase

import (
	"context"

	"github.com/Higor-ViniciusDev/auth/internal/entity"
	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
)

type UserUseCase struct {
	userRepo entity.UserRepository
}

type UserInputValidationDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserOutputValidationDTO struct {
	Token string `json:"token"`
}

func NewUserUseCase(userRepo entity.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (u *UserUseCase) Create(ctx context.Context, user entity.User) *internal_error.InternalError {
	return u.userRepo.Create(ctx, user)
}

func (u *UserUseCase) Delete(ctx context.Context, id string) *internal_error.InternalError {
	return u.userRepo.Delete(ctx, id)
}

func (u *UserUseCase) Validation(ctx context.Context, dto UserInputValidationDTO) *internal_error.InternalError {
	return u.userRepo.Validation(ctx, dto.Email, dto.Password)
}

func (u *UserUseCase) List(ctx context.Context) ([]entity.User, *internal_error.InternalError) {
	return u.userRepo.List(ctx)
}
