package service

import (
	"time"

	"github.com/devWaylander/coins_store/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

type Repository interface {
}

type service struct {
	repo   Repository
	jwtKey string
}

func New(repo Repository, jwtKey string) *service {
	return &service{
		repo:   repo,
		jwtKey: jwtKey,
	}
}

func (s *service) GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
