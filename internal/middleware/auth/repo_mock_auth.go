package auth

import (
	"context"

	"github.com/devWaylander/coins_store/pkg/models"
)

type MockRepository struct {
	CreateUserTXFunc              func(ctx context.Context, username, passwordHash string) (int64, error)
	GetUserByUsernameFunc         func(ctx context.Context, username string) (*models.User, error)
	GetUserPassHashByUsernameFunc func(ctx context.Context, username string) (string, error)
}

func (m *MockRepository) CreateUserTX(ctx context.Context, username, passwordHash string) (int64, error) {
	return m.CreateUserTXFunc(ctx, username, passwordHash)
}

func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return m.GetUserByUsernameFunc(ctx, username)
}

func (m *MockRepository) GetUserPassHashByUsername(ctx context.Context, username string) (string, error) {
	return m.GetUserPassHashByUsernameFunc(ctx, username)
}
