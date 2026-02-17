package entity

import (
	"os"
	"time"

	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
)

type Token struct {
	tokenAuth *jwtauth.JWTAuth
}

type TokenBody struct {
	Subject   string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func NewToken() *Token {
	secret := os.Getenv("JWT_SECRET")

	return &Token{
		tokenAuth: jwtauth.New("HS256", []byte(secret), nil),
	}
}

func NewTokenBody(userID uuid.UUID) *TokenBody {
	now := time.Now()

	return &TokenBody{
		Subject:   userID.String(),
		IssuedAt:  now,
		ExpiresAt: now.Add(time.Minute * 15),
	}
}

func (t *Token) GenerateToken(body *TokenBody) (string, *internal_error.InternalError) {
	_, tokenString, err := t.tokenAuth.Encode(map[string]interface{}{
		"sub": body.Subject,
		"iat": body.IssuedAt.Unix(),
		"exp": body.ExpiresAt.Unix(),
	})

	if err != nil {
		return "", internal_error.NewInternalServerError("error generating token")
	}

	return tokenString, nil
}

func (t *Token) ValidateToken(tokenString string) (*TokenBody, *internal_error.InternalError) {
	token, err := t.tokenAuth.Decode(tokenString)
	if err != nil {
		return nil, internal_error.NewUnauthorizedAccess("invalid token")
	}

	claims := token.PrivateClaims()

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, internal_error.NewUnauthorizedAccess("invalid sub claim")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, internal_error.NewUnauthorizedAccess("invalid exp claim")
	}

	exp := time.Unix(int64(expFloat), 0)

	if time.Now().After(exp) {
		return nil, internal_error.NewUnauthorizedAccess("token expired")
	}

	return &TokenBody{
		Subject:   sub,
		ExpiresAt: exp,
	}, nil
}

type TokenInterface interface {
	GenerateToken(body *TokenBody) (string, *internal_error.InternalError)
	ValidateToken(tokenString string) (*TokenBody, *internal_error.InternalError)
}
