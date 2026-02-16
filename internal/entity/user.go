package entity

import (
	"context"

	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
	"github.com/google/uuid"
)

type User struct {
	id       uuid.UUID
	name     string
	email    string
	password string
}

func NewUser() *User {
	user := &User{
		id: uuid.New(),
	}

	return user
}

func (u *User) GetID() uuid.UUID {
	return u.id
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) SetID(id uuid.UUID) {
	u.id = id
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) SetPassword(password string) {
	u.password = password
}

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *User) *internal_error.InternalError
	Delete(ctx context.Context, id string) *internal_error.InternalError
	Validation(ctx context.Context, email string, password string) *internal_error.InternalError
	List(ctx context.Context) ([]User, *internal_error.InternalError)
	ValidationEmailAlreadyExists(ctx context.Context, email string) *internal_error.InternalError
}
