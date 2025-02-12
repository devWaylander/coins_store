package models

import "github.com/go-openapi/strfmt"

type BalanceDB struct {
	ID        int64            `db:"id"`
	Amount    int64            `db:"amount"`
	DeletedAt *strfmt.DateTime `db:"deleted_at"`
	CreatedAt strfmt.DateTime  `db:"created_at"`
}

func (bdb *BalanceDB) ToModelBalance() *Balance {
	return &Balance{
		ID:     bdb.ID,
		Amount: bdb.Amount,
	}
}

type Balance struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"amount"`
}

func (b *Balance) ToModelBalanceDTO() *BalanceDTO {
	return &BalanceDTO{
		Amount: b.Amount,
	}
}

type BalanceDTO struct {
	Amount int64 `json:"amount"`
}
