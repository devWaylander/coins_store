package models

import "github.com/go-openapi/strfmt"

type BalanceHistoryDB struct {
	ID                int64            `db:"id"`
	BalanceID         int64            `db:"balance_id"`
	TransactionAmount int64            `db:"transaction_amount"`
	Sender            string           `db:"sender"`
	Recipient         string           `db:"recipient"`
	DeletedAt         *strfmt.DateTime `db:"deleted_at"`
	CreatedAt         strfmt.DateTime  `db:"created_at"`
}

func (bhdb *BalanceHistoryDB) ToModelBalanceHistory() BalanceHistory {
	return BalanceHistory{
		TransactionAmount: bhdb.TransactionAmount,
		Sender:            bhdb.Sender,
		Recipient:         bhdb.Recipient,
	}
}

type BalanceHistory struct {
	TransactionAmount int64  `json:"transaction_amount"`
	Sender            string `json:"sender"`
	Recipient         string `json:"recipient"`
}

type ReceivedDTO struct {
	FromUser string `json:"fromUser"`
	Amount   int64  `json:"amount"`
}

type SentDTO struct {
	ToUser string `json:"toUser"`
	Amount int64  `json:"amount"`
}

type BalanceHistoryDTO struct {
	Received []ReceivedDTO `json:"received"`
	Sent     []SentDTO     `json:"sent"`
}
