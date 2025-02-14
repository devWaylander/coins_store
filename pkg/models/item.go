package models

type ItemQuery struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Item     string `json:"item"`
}
