package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/Higor-ViniciusDev/auth/internal/entity"
	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

type UserPG struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Verified     bool      `json:"verified"`
	VerifiedAt   time.Time `json:"verified_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) *internal_error.InternalError {
	hash := sha256.Sum256([]byte(user.GetPassword()))
	hashString := hex.EncodeToString(hash[:])

	query := `INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, user.GetID().String(), user.GetName(), user.GetEmail(), hashString)
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

func (r *UserRepository) Validation(ctx context.Context, email string, password string) (*entity.User, *internal_error.InternalError) {
	hash := sha256.Sum256([]byte(password))
	hashString := hex.EncodeToString(hash[:])

	query := `SELECT id,verified,name,email,password_hash,verified_at,created_at FROM users WHERE email = $1 AND password_hash = $2`
	row := r.db.QueryRowContext(ctx, query, email, hashString)

	var userPG UserPG
	if err := row.Scan(&userPG.ID, &userPG.Verified, &userPG.Name, &userPG.Email, &userPG.PasswordHash, &userPG.VerifiedAt, &userPG.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, internal_error.NewUnauthorizedAccess("Email or password incorrect")
		}
		logger.Error("Error validating user: ", err)
		return nil, internal_error.NewInternalServerError("Error validating user")
	}

	if !userPG.Verified {
		return nil, internal_error.NewUnauthorizedEmailNotVerified("Email not verified")
	}

	var userRetorno = &entity.User{}
	userRetorno.SetID(uuid.MustParse(userPG.ID))
	userRetorno.SetName(userPG.Name)
	userRetorno.SetEmail(userPG.Email)
	userRetorno.SetPassword(userPG.PasswordHash)
	userRetorno.SetVerified(userPG.Verified)
	return userRetorno, nil
}

func (r *UserRepository) List(ctx context.Context) ([]entity.User, *internal_error.InternalError) {
	query := `SELECT id, name, email FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.Error("Error listing users: ", err)
		return nil, internal_error.NewInternalServerError("Error listing users")
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var id uuid.UUID
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			logger.Error("Error scanning user: ", err)
			continue
		}
		user := entity.NewUser()
		user.SetID(id)
		user.SetName(name)
		user.SetEmail(email)
		users = append(users, *user)
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

func (r *UserRepository) MarkVerified(ctx context.Context, id uuid.UUID) (string, *internal_error.InternalError) {
	query := `UPDATE users SET verified = true, verified_at = NOW() WHERE id = $1 AND verified = false RETURNING email`
	row := r.db.QueryRowContext(ctx, query, id)

	var email string
	if err := row.Scan(&email); err != nil {
		if err == sql.ErrNoRows {
			return "", internal_error.NewBadRequestError("email already verified or user not found")
		}
		logger.Error("Error marking user as verified: ", err)
		return "", internal_error.NewInternalServerError("Error marking user as verified")
	}

	return email, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, *internal_error.InternalError) {
	query := `SELECT id, name, email FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var id uuid.UUID
	var name, userEmail string
	if err := row.Scan(&id, &name, &userEmail); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // not found — not an error (security)
		}
		logger.Error("Error finding user by email: ", err)
		return nil, internal_error.NewInternalServerError("Error finding user")
	}

	user := entity.NewUser()
	user.SetID(id)
	user.SetName(name)
	user.SetEmail(userEmail)
	return user, nil
}
