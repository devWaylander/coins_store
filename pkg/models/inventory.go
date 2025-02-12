package models

import "github.com/go-openapi/strfmt"

type InventoryDB struct {
	ID        int64            `db:"id"`
	MerchID   int64            `db:"merch_id"`
	Count     int64            `db:"count"`
	DeletedAt *strfmt.DateTime `db:"deleted_at"`
	CreatedAt strfmt.DateTime  `db:"created_at"`
}

func (idb *InventoryDB) ToModelInventory() *Inventory {
	return &Inventory{
		ID:      idb.ID,
		MerchID: idb.MerchID,
		Count:   idb.Count,
	}
}

type Inventory struct {
	ID      int64 `json:"id"`
	MerchID int64 `json:"merch_id"`
	Count   int64 `json:"count"`
}

func (i *Inventory) ToModelInventoryDTO() *InventoryDTO {
	return &InventoryDTO{
		MerchID: i.MerchID,
		Count:   i.Count,
	}
}

type InventoryDTO struct {
	MerchID int64 `json:"merch_id"`
	Count   int64 `json:"count"`
}
