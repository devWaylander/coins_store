package models

import "github.com/go-openapi/strfmt"

type MerchDB struct {
	ID        int64            `db:"id"`
	Name      int64            `db:"name"`
	Price     int64            `db:"price"`
	DeletedAt *strfmt.DateTime `db:"deleted_at"`
	CreatedAt strfmt.DateTime  `db:"created_at"`
}

func (mdb *MerchDB) ToModelMerch() *Merch {
	return &Merch{
		Name:  mdb.Name,
		Price: mdb.Price,
	}
}

type Merch struct {
	ID    int64 `json:"id"`
	Name  int64 `json:"name"`
	Price int64 `json:"price"`
}
