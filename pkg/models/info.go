package models

type InfoDTO struct {
	Coins        int64               `json:"coins"`
	Inventory    []MerchDTO          `json:"inventory"`
	CoinsHistory []BalanceHistoryDTO `json:"coinHistory"`
}
