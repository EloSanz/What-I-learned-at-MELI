package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

// Service defines the business logic for authentication
type Service interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type service struct {
	secretKey []byte
}

// NewService creates a new auth service instance
func NewService() Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// In a real application, failing fast here is best practice.
		// We'll panic to ensure the service doesn't start without a secret.
		panic("JWT_SECRET environment variable is strictly required")
	}

	return &service{
		secretKey: []byte(secret),
	}
}

// GenerateToken creates a new JWT with a 1-hour expiration
func (s *service) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(s.secretKey)
}

// ValidateToken parses and validates the JWT
func (s *service) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil
}
