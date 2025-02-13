package service

import (
	"context"
	"errors"

	internalErrors "github.com/devWaylander/coins_store/pkg/errors"
	"github.com/devWaylander/coins_store/pkg/models"
)

type Repository interface {
	// Balance
	GetBalanceByUserID(ctx context.Context, userID int64) (models.Balance, error)
	GetBalanceAmountByUserID(ctx context.Context, userID int64) (int64, error)
	// Balance history
	GetBalanceHistoryByUserID(ctx context.Context, userID int64) ([]models.BalanceHistory, error)
	// Inventory
	GetInventoryMerchItems(ctx context.Context, userID int64) ([]models.InventoryMerch, error)
	GetInventoryIDByUserID(ctx context.Context, userID int64) (int64, error)
	// Merch
	GetMerchByName(ctx context.Context, name string) (models.Merch, error)
	BuyItemTX(ctx context.Context, userID, balanceID, inventoryID, merchID, price int64, username, item string) error
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

func (s *service) GetUserInfo(ctx context.Context, userID int64, username string) (models.InfoDTO, error) {
	info := models.InfoDTO{}

	// Balance
	amount, err := s.getBalanceAmount(ctx, userID)
	if err != nil {
		return models.InfoDTO{}, err
	}
	info.Coins = amount

	// CoinsHistory
	balanceHistory, err := s.getBalanceHistory(ctx, userID, username)
	if err != nil {
		return models.InfoDTO{}, err
	}
	info.CoinsHistory = balanceHistory

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

	return info, nil
}

func (s *service) getBalanceAmount(ctx context.Context, userID int64) (int64, error) {
	amount, err := s.repo.GetBalanceAmountByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

func (s *service) getBalanceHistory(ctx context.Context, userID int64, username string) (models.BalanceHistoryDTO, error) {
	balanceHistory, err := s.repo.GetBalanceHistoryByUserID(ctx, userID)
	if err != nil {
		return models.BalanceHistoryDTO{}, err
	}

	var received = []models.ReceivedDTO{}
	var sent = []models.SentDTO{}
	for _, item := range balanceHistory {
		if item.Recipient == username {
			received = append(received, models.ReceivedDTO{
				FromUser: item.Sender,
				Amount:   item.TransactionAmount,
			})
			continue
		}

		sent = append(sent, models.SentDTO{
			ToUser: item.Recipient,
			Amount: item.TransactionAmount,
		})
	}

	return models.BalanceHistoryDTO{Received: received, Sent: sent}, nil
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

func (s *service) BuyItem(ctx context.Context, userID int64, username, item string) error {
	merch, err := s.repo.GetMerchByName(ctx, item)
	if err != nil {
		return err
	}
	if merch.ID == 0 {
		return errors.New(internalErrors.ErrItemDoesntExist)
	}

	balance, err := s.repo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if balance.Amount-merch.Price < 0 {
		return errors.New(internalErrors.ErrNotEnoughCoins)
	}

	inventoryID, err := s.repo.GetInventoryIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	err = s.repo.BuyItemTX(ctx, userID, balance.ID, inventoryID, merch.ID, merch.Price, username, merch.Name)
	if err != nil {
		return err
	}

	return nil
}
