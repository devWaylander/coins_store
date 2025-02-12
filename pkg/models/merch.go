package models

import "github.com/go-openapi/strfmt"

type MerchDB struct {
	ID        int64            `db:"id"`
	Name      string           `db:"name"`
	Price     int64            `db:"price"`
	DeletedAt *strfmt.DateTime `db:"deleted_at"`
	CreatedAt strfmt.DateTime  `db:"created_at"`
}

func (mdb *MerchDB) ToModelMerch() Merch {
	return Merch{
		Name:  mdb.Name,
		Price: mdb.Price,
	}
}

type Merch struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Count int64  `json:"count"`
}

func (m *Merch) ToModelMerchDTO() MerchDTO {
	return MerchDTO{
		Type:     m.Name,
		Quantity: m.Count,
	}
}

type MerchDTO struct {
	Type     string `json:"type"`
	Quantity int64  `json:"quantity"`
}
