package models

type SendCoinsReqBody struct {
	Recipient string `json:"toUser"`
	Amount    int64  `json:"amount"`
}

type CoinsQuery struct {
	UserID    int64  `json:"user_id"`
	Amount    int64  `json:"amount"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}
