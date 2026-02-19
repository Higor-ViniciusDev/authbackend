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
	verified bool
}

func NewUser() *User {
	return &User{id: uuid.New()}
}

func (u *User) GetID() uuid.UUID    { return u.id }
func (u *User) GetName() string     { return u.name }
func (u *User) GetEmail() string    { return u.email }
func (u *User) GetPassword() string { return u.password }
func (u *User) GetVerified() bool   { return u.verified }

func (u *User) SetID(id uuid.UUID)        { u.id = id }
func (u *User) SetName(name string)       { u.name = name }
func (u *User) SetEmail(email string)     { u.email = email }
func (u *User) SetPassword(pwd string)    { u.password = pwd }
func (u *User) SetVerified(verified bool) { u.verified = verified }

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *User) *internal_error.InternalError
	Delete(ctx context.Context, id string) *internal_error.InternalError
	Validation(ctx context.Context, email string, password string) (*User, *internal_error.InternalError)
	List(ctx context.Context) ([]User, *internal_error.InternalError)
	ValidationEmailAlreadyExists(ctx context.Context, email string) *internal_error.InternalError
	MarkVerified(ctx context.Context, id uuid.UUID) (string, *internal_error.InternalError)
	FindByEmail(ctx context.Context, email string) (*User, *internal_error.InternalError)
}
