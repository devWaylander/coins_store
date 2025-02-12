package service

import (
	"context"

	"github.com/devWaylander/coins_store/pkg/models"
)

type Repository interface {
	GetBalanceByUserID(ctx context.Context, userID int64) (int64, error)
	GetInventoryMerchItems(ctx context.Context, userID int64) ([]models.InventoryMerch, error)
	// GetInventoryMerchesByIDs(ctx context.Context, merchesIDs []int64) ([]models.Merch, error)
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
	inventory, err := s.getInventory(ctx, userID)
	if err != nil {
		return models.InfoDTO{}, err
	}
	itemsDTO := make([]models.MerchDTO, 0, len(inventory.Items))
	for _, item := range inventory.Items {
		itemsDTO = append(itemsDTO, item.ToModelMerchDTO())
	}
	info.Inventory = itemsDTO

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

func (s *service) getInventory(ctx context.Context, userID int64) (models.Inventory, error) {
	inventoryMerchItems, err := s.repo.GetInventoryMerchItems(ctx, userID)
	if err != nil {
		return models.Inventory{}, err
	}

	inventory := models.Inventory{Items: make([]models.Merch, 0, len(inventoryMerchItems))}
	for _, item := range inventoryMerchItems {
		inventory.Items = append(inventory.Items, models.Merch{
			ID:    item.MerchID,
			Name:  item.Name,
			Count: item.Count,
		})
	}

	return inventory, nil
}
