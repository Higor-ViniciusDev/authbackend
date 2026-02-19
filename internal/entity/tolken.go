package entity

import (
	"fmt"
	"os"
	"time"

	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Token struct {
	secret []byte
}

type TokenBody struct {
	Subject   string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func NewToken() *Token {
	secret := os.Getenv("JWT_SECRET")

	return &Token{
		secret: []byte(secret),
	}
}

func NewTokenBody(userID uuid.UUID) *TokenBody {
	now := time.Now()

	return &TokenBody{
		Subject:   userID.String(),
		IssuedAt:  now,
		ExpiresAt: now.Add(time.Hour * 12),
	}
}

func (t *Token) GenerateToken(body *TokenBody) (string, *internal_error.InternalError) {
	claims := jwt.MapClaims{
		"sub": body.Subject,
		"iat": body.IssuedAt.Unix(),
		"exp": body.ExpiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(t.secret)
	if err != nil {
		return "", internal_error.NewInternalServerError("error generating token")
	}

	return tokenString, nil
}

// ValidateToken decodes and validates the JWT.
// returns error if: invalid token, expired, or bad claims.
func (t *Token) ValidateToken(tokenString string) (*TokenBody, *internal_error.InternalError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// protects against algorithm switching
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, internal_error.NewUnauthorizedAccess("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, internal_error.NewUnauthorizedAccess("invalid claims")
	}

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

	var iat time.Time
	if iatFloat, ok := claims["iat"].(float64); ok {
		iat = time.Unix(int64(iatFloat), 0)
	}

	return &TokenBody{
		Subject:   sub,
		IssuedAt:  iat,
		ExpiresAt: exp,
	}, nil
}

type TokenInterface interface {
	GenerateToken(body *TokenBody) (string, *internal_error.InternalError)
	ValidateToken(tokenString string) (*TokenBody, *internal_error.InternalError)
}
