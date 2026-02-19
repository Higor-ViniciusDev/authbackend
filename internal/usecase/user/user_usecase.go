package user_usecase

import (
	"context"

	"github.com/Higor-ViniciusDev/auth/internal/entity"
	events_internal "github.com/Higor-ViniciusDev/auth/internal/infra/events"
	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
	"github.com/Higor-ViniciusDev/auth/pkg/events"
	"github.com/google/uuid"
)

type UserUseCase struct {
	userRepo        entity.UserRepositoryInterface
	EventDispatched events.EventDispachtInterface
	Token           entity.TokenInterface
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
	Token   string `json:"token"`
	Pedding bool   `json:"pedding"`
}

type UserSendEmailConfirmationDTO struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ResendVerificationDTO struct {
	Email string `json:"email"`
}

type UserEmailVerifiedDTO struct {
	Email string `json:"email"`
}

func NewUserUseCase(userRepo entity.UserRepositoryInterface, DisparadorEvento events.EventDispachtInterface, token entity.TokenInterface) *UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		EventDispatched: DisparadorEvento,
		Token:           token,
	}
}

func (u *UserUseCase) Create(ctx context.Context, dto CreateUserDTO) *internal_error.InternalError {
	user := entity.NewUser()
	user.SetName(dto.Name)
	user.SetEmail(dto.Email)
	user.SetPassword(dto.Password)

	if err := u.userRepo.ValidationEmailAlreadyExists(ctx, user.GetEmail()); err != nil {
		return err
	}

	retorno := u.userRepo.Create(ctx, user)

	if retorno != nil {
		return retorno
	}

	return u.dispatchVerificationEmail(user.GetID(), user.GetEmail())
}

// VerifyEmail validates the JWT from the link, updates verified=true in the db
// and fires UserEmailVerified to notify all pods via fanout.
func (u *UserUseCase) VerifyEmail(ctx context.Context, tokenString string) *internal_error.InternalError {
	// decodes and validates the JWT (exp + purpose via sub)
	tokenBody, err := u.Token.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	userID, parseErr := uuid.Parse(tokenBody.Subject)
	if parseErr != nil {
		return internal_error.NewUnauthorizedAccess("token invalid: sub not is UUID")
	}

	// sets verified=true — AND verified=false prevents double processing
	email, err := u.userRepo.MarkVerified(ctx, userID)
	if err != nil {
		return err
	}

	// publishes to fanout exchange → all pods receive and notify WS if they have the connection
	verifiedEvent := events_internal.NewUserEmailVerified()
	verifiedEvent.SetPayload(UserEmailVerifiedDTO{Email: email})
	u.EventDispatched.Dispatch(verifiedEvent)

	return nil
}

// ResendVerification finds user by email, generates new JWT and publishes to email.pending queue.
func (u *UserUseCase) ResendVerification(ctx context.Context, dto ResendVerificationDTO) *internal_error.InternalError {
	user, err := u.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		// always returns OK for security (doesn't reveal if email exists)
		return nil
	}

	if user == nil {
		return nil
	}

	return u.dispatchVerificationEmail(user.GetID(), user.GetEmail())
}

func (u *UserUseCase) Delete(ctx context.Context, id string) *internal_error.InternalError {
	return u.userRepo.Delete(ctx, id)
}

func (u *UserUseCase) Validation(ctx context.Context, dto UserInputValidationDTO) (*UserOutputValidationDTO, *internal_error.InternalError) {
	retorno, err := u.userRepo.Validation(ctx, dto.Email, dto.Password)
	if err != nil {
		return nil, err
	}

	tokenString, err := u.Token.GenerateToken(entity.NewTokenBody(retorno.GetID()))
	if err != nil {
		return nil, err
	}

	return &UserOutputValidationDTO{
		Token:   tokenString,
		Pedding: retorno.GetVerified(),
	}, nil
}

func (u *UserUseCase) List(ctx context.Context) ([]entity.User, *internal_error.InternalError) {
	return u.userRepo.List(ctx)
}

func (u *UserUseCase) dispatchVerificationEmail(userID uuid.UUID, email string) *internal_error.InternalError {
	tokenBody := entity.NewTokenBody(userID)
	tokenString, err := u.Token.GenerateToken(tokenBody)
	if err != nil {
		return err
	}

	event := events_internal.NewUserCreated()
	event.SetPayload(UserSendEmailConfirmationDTO{
		Email: email,
		Token: tokenString,
	})

	u.EventDispatched.Dispatch(event)
	return nil
}
