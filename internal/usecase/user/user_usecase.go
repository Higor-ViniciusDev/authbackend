package user_usecase

import (
	"context"

	"github.com/Higor-ViniciusDev/auth/internal/entity"
	events_internal "github.com/Higor-ViniciusDev/auth/internal/infra/events"
	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
	"github.com/Higor-ViniciusDev/auth/pkg/events"
)

type UserUseCase struct {
	userRepo        entity.UserRepositoryInterface
	EventDispatched events.EventDispachtInterface
}

type CreateUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInputValidationDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserOutputValidationDTO struct {
	Token string `json:"token"`
}

func NewUserUseCase(userRepo entity.UserRepositoryInterface, DisparadorEvento events.EventDispachtInterface) *UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		EventDispatched: DisparadorEvento,
	}
}

func (u *UserUseCase) Create(ctx context.Context, dto CreateUserDTO) *internal_error.InternalError {
	user := entity.NewUser()
	user.SetName(dto.Name)
	user.SetEmail(dto.Email)
	user.SetPassword(dto.Password)

	retorno := u.userRepo.Create(ctx, user)

	if retorno != nil {
		return retorno
	}

	newUserCreatedEvent := events_internal.NewUserCreated()
	newUserCreatedEvent.SetPayload(dto.Email)
	u.EventDispatched.Dispatch(newUserCreatedEvent)
	return nil
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
