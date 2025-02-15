package service

import (
	"context"

	"github.com/devWaylander/coins_store/pkg/models"
)

type MockRepository struct {
	IsUserExistFunc               func(ctx context.Context, username string) (bool, error)
	GetBalanceIDByUsernameFunc    func(ctx context.Context, username string) (int64, error)
	GetBalanceByUserIDFunc        func(ctx context.Context, userID int64) (models.Balance, error)
	GetBalanceAmountByUserIDFunc  func(ctx context.Context, userID int64) (int64, error)
	GetBalanceHistoryByUserIDFunc func(ctx context.Context, userID int64) ([]models.BalanceHistory, error)
	GetInventoryMerchItemsFunc    func(ctx context.Context, userID int64) ([]models.InventoryMerch, error)
	GetInventoryIDByUserIDFunc    func(ctx context.Context, userID int64) (int64, error)
	GetMerchByNameFunc            func(ctx context.Context, name string) (models.Merch, error)
	BuyItemTXFunc                 func(ctx context.Context, userID, balanceID, inventoryID, merchID, price int64, username, item string) error
	SendCoinsTXFunc               func(ctx context.Context, userID, senderBalanceID, recipientBalanceID, amount int64, sender, recipient string) error
}

func (m *MockRepository) IsUserExist(ctx context.Context, username string) (bool, error) {
	return m.IsUserExistFunc(ctx, username)
}

func (m *MockRepository) GetBalanceIDByUsername(ctx context.Context, username string) (int64, error) {
	return m.GetBalanceIDByUsernameFunc(ctx, username)
}

func (m *MockRepository) GetBalanceByUserID(ctx context.Context, userID int64) (models.Balance, error) {
	return m.GetBalanceByUserIDFunc(ctx, userID)
}

func (m *MockRepository) GetBalanceAmountByUserID(ctx context.Context, userID int64) (int64, error) {
	return m.GetBalanceAmountByUserIDFunc(ctx, userID)
}

func (m *MockRepository) GetBalanceHistoryByUserID(ctx context.Context, userID int64) ([]models.BalanceHistory, error) {
	return m.GetBalanceHistoryByUserIDFunc(ctx, userID)
}

func (m *MockRepository) GetInventoryMerchItems(ctx context.Context, userID int64) ([]models.InventoryMerch, error) {
	return m.GetInventoryMerchItemsFunc(ctx, userID)
}

func (m *MockRepository) GetInventoryIDByUserID(ctx context.Context, userID int64) (int64, error) {
	return m.GetInventoryIDByUserIDFunc(ctx, userID)
}

func (m *MockRepository) GetMerchByName(ctx context.Context, name string) (models.Merch, error) {
	return m.GetMerchByNameFunc(ctx, name)
}

func (m *MockRepository) BuyItemTX(ctx context.Context, userID, balanceID, inventoryID, merchID, price int64, username, item string) error {
	return m.BuyItemTXFunc(ctx, userID, balanceID, inventoryID, merchID, price, username, item)
}

func (m *MockRepository) SendCoinsTX(ctx context.Context, userID, senderBalanceID, recipientBalanceID, amount int64, sender, recipient string) error {
	return m.SendCoinsTXFunc(ctx, userID, senderBalanceID, recipientBalanceID, amount, sender, recipient)
}
