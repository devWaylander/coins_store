package models

import "github.com/go-openapi/strfmt"

type InventoryDB struct {
	ID        int64            `db:"id"`
	UserID    int64            `db:"user_id"`
	DeletedAt *strfmt.DateTime `db:"deleted_at"`
	CreatedAt strfmt.DateTime  `db:"created_at"`
}

type Inventory struct {
	Items []Merch `json:"items"`
}

type InventoryMerchDB struct {
	InventoryID int64            `db:"inventory_id"`
	MerchID     int64            `db:"merch_id"`
	Name        string           `db:"name"`
	Count       int64            `db:"count"`
	DeletedAt   *strfmt.DateTime `db:"deleted_at"`
	CreatedAt   strfmt.DateTime  `db:"created_at"`
}

func (imdb *InventoryMerchDB) ToModelInventoryMerch() InventoryMerch {
	return InventoryMerch{
		InventoryID: imdb.InventoryID,
		MerchID:     imdb.MerchID,
		Name:        imdb.Name,
		Count:       imdb.Count,
	}
}

type InventoryMerch struct {
	InventoryID int64  `json:"inventory_id"`
	MerchID     int64  `json:"merch_id"`
	Name        string `json:"name"`
	Count       int64  `json:"count"`
}
