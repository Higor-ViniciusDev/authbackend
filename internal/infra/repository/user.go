package repository

import (
	"context"

	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/Higor-ViniciusDev/auth/internal/entity"
	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
	"github.com/go-pg/pg/v10"
)

type UserRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) *internal_error.InternalError {
	_, err := r.db.Model(&user).Insert()
	if err != nil {
		logger.Error("Error create user in Repository: ", err)
		return internal_error.NewInternalServerError("Error to create user")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) *internal_error.InternalError {
	_, err := r.db.Model(&entity.User{}).Where("id = ?", id).Delete()
	if err != nil {
		logger.Error("Error delete user in Repository: ", err)
		return internal_error.NewInternalServerError("Error to delete user")
	}
	return nil
}

func (r *UserRepository) Validation(ctx context.Context, email string, password string) *internal_error.InternalError {
	var user entity.User
	err := r.db.Model(&user).Where("email = ?", email).Where("password = ?", password).Select()
	if err != nil {
		logger.Error("Error validation user in Repository: ", err)
		return internal_error.NewUnauthorizedAccess("Email ou login incorreto")
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]entity.User, *internal_error.InternalError) {
	var users []entity.User
	err := r.db.Model(&users).Select()

	if err != nil {
		logger.Error("Error list user in Repository: ", err)
		return nil, internal_error.NewInternalServerError("Error to list user")
	}

	return users, nil
}
