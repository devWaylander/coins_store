package service

import (
	"context"

	"github.com/devWaylander/coins_store/pkg/models"
)

type Repository interface {
	GetBalanceByUserID(ctx context.Context, userID int64) (int64, error)
}

type service struct {
	repo Repository
}

func New(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetUserInfo(ctx context.Context, userID int64) (models.InfoDTO, error) {
	info := models.InfoDTO{}
	// Balance
	balance, err := s.getBalance(ctx, userID)
	if err != nil {
		return models.InfoDTO{}, err
	}
	info.Coins = balance

	// Inventory

	// CoinsHistory

	return info, nil
}

func (s *service) getBalance(ctx context.Context, userID int64) (int64, error) {
	balance, err := s.repo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
