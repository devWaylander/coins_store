package models

import "github.com/go-openapi/strfmt"

type UserDB struct {
	ID           int64            `db:"id"`
	BalanceID    int64            `db:"balance_id"`
	InventoryID  *int64           `db:"inventory_id"`
	Username     string           `db:"username"`
	PasswordHash string           `db:"password_hash"`
	DeletedAt    *strfmt.DateTime `db:"deleted_at"`
	CreatedAt    strfmt.DateTime  `db:"created_at"`
}

func (udb *UserDB) ToModelUser() *User {
	return &User{
		ID:           udb.ID,
		BalanceID:    udb.BalanceID,
		InventoryID:  udb.InventoryID,
		Username:     udb.Username,
		PasswordHash: udb.PasswordHash,
	}
}

type User struct {
	ID           int64  `json:"id"`
	BalanceID    int64  `json:"balance_id"`
	InventoryID  *int64 `json:"inventory_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}
