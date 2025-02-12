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
	Count       int64            `db:"count"`
	DeletedAt   *strfmt.DateTime `db:"deleted_at"`
	CreatedAt   strfmt.DateTime  `db:"created_at"`
}

func (imdb *InventoryMerchDB) ToModelInventory(inventoryItems []InventoryMerchDB, merchMap map[int64]*MerchDB) *Inventory {
	inventory := Inventory{}

	for _, item := range inventoryItems {
		if merch, exists := merchMap[item.MerchID]; exists {
			inventory.Items = append(inventory.Items, Merch{
				Name:  merch.Name,
				Price: merch.Price,
				Count: item.Count,
			})
		}
	}

	return &inventory
}
