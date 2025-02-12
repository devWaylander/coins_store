package service

import (
	"context"
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

// Auth
func (s *service) generateJWT(userID int64) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.Claims{
		UserID: userID,
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

// User
func (s *service) LoginWithPass(ctx context.Context, username, password string) error {
	// гетаемюзера, если нет, то сразу регаем (создаём юзернейм, пароль, инвентарь)
	// если есть, то валидируем его и отдаём токен

	return nil
}
