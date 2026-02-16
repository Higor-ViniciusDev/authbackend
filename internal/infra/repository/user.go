package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/Higor-ViniciusDev/auth/internal/entity"
	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) *internal_error.InternalError {

	if err := r.ValidationEmailAlreadyExists(ctx, user.GetEmail()); err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)`

	//convert senha 256
	hash := sha256.Sum256([]byte(user.GetPassword()))
	hashString := hex.EncodeToString(hash[:])

	_, err := r.db.ExecContext(ctx, query, user.GetName(), user.GetEmail(), hashString)
	if err != nil {
		logger.Error("Error creating user: ", err)
		return internal_error.NewInternalServerError("Error creating user")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) *internal_error.InternalError {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("Error deleting user: ", err)
		return internal_error.NewInternalServerError("Error deleting user")
	}
	return nil
}

func (r *UserRepository) Validation(ctx context.Context, email string, password string) *internal_error.InternalError {
	query := `SELECT id, name, email, password FROM users WHERE email = $1 AND password = $2`

	row := r.db.QueryRowContext(ctx, query, email, password)

	var id uuid.UUID
	var name, userEmail, userPassword string

	err := row.Scan(&id, &name, &userEmail, &userPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return internal_error.NewUnauthorizedAccess("Email or password incorrect")
		}
		logger.Error("Error validating user: ", err)
		return internal_error.NewInternalServerError("Error validating user")
	}

	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]entity.User, *internal_error.InternalError) {
	query := `SELECT id, name, email, password FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.Error("Error listing users: ", err)
		return nil, internal_error.NewInternalServerError("Error listing users")
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var id uuid.UUID
		var name, email, password string

		if err := rows.Scan(&id, &name, &email, &password); err != nil {
			logger.Error("Error scanning user: ", err)
			continue
		}

		user := entity.NewUser()
		user.SetID(id)
		user.SetName(name)
		user.SetEmail(email)
		user.SetPassword(password)

		users = append(users, *user)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Error iterating users: ", err)
		return nil, internal_error.NewInternalServerError("Error iterating users")
	}

	return users, nil
}

func (r *UserRepository) ValidationEmailAlreadyExists(ctx context.Context, email string) *internal_error.InternalError {
	query := `SELECT id FROM users WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)
	var id string
	_ = row.Scan(&id)
	if id != "" {
		return internal_error.NewUnauthorizedEmailAlreadyExists("Email already exists")
	}
	return nil
}
