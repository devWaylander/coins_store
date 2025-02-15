package models

type InfoQuery struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

type InfoDTO struct {
	Coins        int64             `json:"coins"`
	Inventory    []MerchDTO        `json:"inventory"`
	CoinsHistory BalanceHistoryDTO `json:"coinHistory"`
}
