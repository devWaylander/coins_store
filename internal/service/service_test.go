package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/devWaylander/coins_store/pkg/models"
)

func Test_service_GetUserInfo(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		ctx context.Context
		qp  models.InfoQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.InfoDTO
		wantErr bool
	}{
		{
			name: "success_-_user_exists",
			fields: fields{
				repo: &MockRepository{
					GetBalanceAmountByUserIDFunc: func(ctx context.Context, userID int64) (int64, error) {
						return 200, nil
					},
					GetBalanceHistoryByUserIDFunc: func(ctx context.Context, userID int64) ([]models.BalanceHistory, error) {
						return []models.BalanceHistory{
							{TransactionAmount: 100, Sender: "user2", Recipient: "user1"},
							{TransactionAmount: 50, Sender: "user1", Recipient: "user2"},
						}, nil
					},
					GetInventoryMerchItemsFunc: func(ctx context.Context, userID int64) ([]models.InventoryMerch, error) {
						return []models.InventoryMerch{
							{InventoryID: 1, MerchID: 201, Name: "t-shirt", Count: 15},
						}, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.InfoQuery{Username: "user1"},
			},
			want: models.InfoDTO{
				Coins: 200,
				Inventory: []models.MerchDTO{
					{Type: "t-shirt", Quantity: 15},
				},
				CoinsHistory: models.BalanceHistoryDTO{
					Received: []models.ReceivedDTO{
						{FromUser: "user2", Amount: 100},
					},
					Sent: []models.SentDTO{
						{ToUser: "user2", Amount: 50},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error_-_repository_returns_GetBalanceAmountByUserIDFunc_error",
			fields: fields{
				repo: &MockRepository{
					GetBalanceAmountByUserIDFunc: func(ctx context.Context, userID int64) (int64, error) {
						return 0, errors.New("fail")
					},
					GetBalanceHistoryByUserIDFunc: func(ctx context.Context, userID int64) ([]models.BalanceHistory, error) {
						return []models.BalanceHistory{
							{TransactionAmount: 100, Sender: "user2", Recipient: "user1"},
							{TransactionAmount: 50, Sender: "user1", Recipient: "user2"},
						}, nil
					},
					GetInventoryMerchItemsFunc: func(ctx context.Context, userID int64) ([]models.InventoryMerch, error) {
						return []models.InventoryMerch{
							{InventoryID: 1, MerchID: 201, Name: "t-shirt", Count: 15},
						}, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.InfoQuery{Username: "user1"},
			},
			want:    models.InfoDTO{},
			wantErr: true,
		},
		{
			name: "error_-_repository_returns_GetBalanceHistoryByUserIDFunc_error",
			fields: fields{
				repo: &MockRepository{
					GetBalanceAmountByUserIDFunc: func(ctx context.Context, userID int64) (int64, error) {
						return 200, nil
					},
					GetBalanceHistoryByUserIDFunc: func(ctx context.Context, userID int64) ([]models.BalanceHistory, error) {
						return []models.BalanceHistory{}, errors.New("fail")
					},
					GetInventoryMerchItemsFunc: func(ctx context.Context, userID int64) ([]models.InventoryMerch, error) {
						return []models.InventoryMerch{
							{InventoryID: 1, MerchID: 201, Name: "t-shirt", Count: 15},
						}, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.InfoQuery{Username: "user1"},
			},
			want:    models.InfoDTO{},
			wantErr: true,
		},
		{
			name: "error_-_repository_returns_GetInventoryMerchItemsFunc_error",
			fields: fields{
				repo: &MockRepository{
					GetBalanceAmountByUserIDFunc: func(ctx context.Context, userID int64) (int64, error) {
						return 200, nil
					},
					GetBalanceHistoryByUserIDFunc: func(ctx context.Context, userID int64) ([]models.BalanceHistory, error) {
						return []models.BalanceHistory{
							{TransactionAmount: 100, Sender: "user2", Recipient: "user1"},
							{TransactionAmount: 50, Sender: "user1", Recipient: "user2"},
						}, nil
					},
					GetInventoryMerchItemsFunc: func(ctx context.Context, userID int64) ([]models.InventoryMerch, error) {
						return []models.InventoryMerch{}, errors.New("fail")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.InfoQuery{Username: "user1"},
			},
			want:    models.InfoDTO{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.fields.repo,
			}
			got, err := s.GetUserInfo(tt.args.ctx, tt.args.qp)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetUserInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_BuyItem(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		ctx context.Context
		qp  models.ItemQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success_-_item_purchased",
			fields: fields{
				repo: &MockRepository{
					GetMerchByNameFunc: func(ctx context.Context, name string) (models.Merch, error) {
						return models.Merch{ID: 1, Name: "t-shirt", Price: 500}, nil
					},
					GetBalanceByUserIDFunc: func(ctx context.Context, userID int64) (models.Balance, error) {
						return models.Balance{ID: 1, Amount: 1000}, nil
					},
					GetInventoryIDByUserIDFunc: func(ctx context.Context, userID int64) (int64, error) {
						return 10, nil
					},
					BuyItemTXFunc: func(ctx context.Context, userID, balanceID, inventoryID, merchID, price int64, username, item string) error {
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.ItemQuery{UserID: 1, Username: "user1", Item: "t-shirt"},
			},
			wantErr: false,
		},
		{
			name: "error_-_item_doesn't_exist",
			fields: fields{
				repo: &MockRepository{
					GetMerchByNameFunc: func(ctx context.Context, name string) (models.Merch, error) {
						return models.Merch{}, errors.New("fail")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.ItemQuery{UserID: 1, Username: "user1", Item: "NonExistentItem"},
			},
			wantErr: true,
		},
		{
			name: "error_-_not_enough_coins",
			fields: fields{
				repo: &MockRepository{
					GetMerchByNameFunc: func(ctx context.Context, name string) (models.Merch, error) {
						return models.Merch{ID: 1, Name: "t-shirt", Price: 1000}, nil
					},
					GetBalanceByUserIDFunc: func(ctx context.Context, userID int64) (models.Balance, error) {
						return models.Balance{ID: 1, Amount: 200}, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.ItemQuery{UserID: 1, Username: "user1", Item: "t-shirt"},
			},
			wantErr: true,
		},
		{
			name: "error_-_database_error_on_balance_retrieval",
			fields: fields{
				repo: &MockRepository{
					GetMerchByNameFunc: func(ctx context.Context, name string) (models.Merch, error) {
						return models.Merch{ID: 1, Name: "t-shirt", Price: 100}, nil
					},
					GetBalanceByUserIDFunc: func(ctx context.Context, userID int64) (models.Balance, error) {
						return models.Balance{}, errors.New("db error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.ItemQuery{UserID: 1, Username: "user1", Item: "t-shirt"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.fields.repo,
			}
			if err := s.BuyItem(tt.args.ctx, tt.args.qp); (err != nil) != tt.wantErr {
				t.Errorf("service.BuyItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_SendCoins(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		ctx context.Context
		qp  models.CoinsQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success_-_send_coins",
			fields: fields{
				repo: &MockRepository{
					IsUserExistFunc: func(ctx context.Context, username string) (bool, error) {
						return true, nil
					},
					GetBalanceByUserIDFunc: func(ctx context.Context, userID int64) (models.Balance, error) {
						return models.Balance{ID: 1, Amount: 1000}, nil
					},
					GetBalanceIDByUsernameFunc: func(ctx context.Context, username string) (int64, error) {
						return 2, nil
					},
					SendCoinsTXFunc: func(ctx context.Context, userID, senderBalanceID, recipientBalanceID, amount int64, sender, recipient string) error {
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.CoinsQuery{UserID: 1, Recipient: "user2", Amount: 500, Sender: "user1"},
			},
			wantErr: false,
		},
		{
			name: "error_-_recipient_does_not_exist",
			fields: fields{
				repo: &MockRepository{
					IsUserExistFunc: func(ctx context.Context, username string) (bool, error) {
						return false, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.CoinsQuery{UserID: 1, Recipient: "user_invalid", Amount: 100, Sender: "user1"},
			},
			wantErr: true,
		},
		{
			name: "error_-_not_enough_coins",
			fields: fields{
				repo: &MockRepository{
					IsUserExistFunc: func(ctx context.Context, username string) (bool, error) {
						return true, nil
					},
					GetBalanceByUserIDFunc: func(ctx context.Context, userID int64) (models.Balance, error) {
						return models.Balance{ID: 1, Amount: 50}, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.CoinsQuery{UserID: 1, Recipient: "user2", Amount: 100, Sender: "user1"},
			},
			wantErr: true,
		},
		{
			name: "error_-_repository_failure",
			fields: fields{
				repo: &MockRepository{
					IsUserExistFunc: func(ctx context.Context, username string) (bool, error) {
						return true, nil
					},
					GetBalanceByUserIDFunc: func(ctx context.Context, userID int64) (models.Balance, error) {
						return models.Balance{ID: 1, Amount: 1000}, nil
					},
					GetBalanceIDByUsernameFunc: func(ctx context.Context, username string) (int64, error) {
						return 2, nil
					},
					SendCoinsTXFunc: func(ctx context.Context, userID, senderBalanceID, recipientBalanceID, amount int64, sender, recipient string) error {
						return errors.New("db error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				qp:  models.CoinsQuery{UserID: 1, Recipient: "user2", Amount: 500, Sender: "user1"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.fields.repo,
			}
			if err := s.SendCoins(tt.args.ctx, tt.args.qp); (err != nil) != tt.wantErr {
				t.Errorf("service.SendCoins() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
