package models

import "github.com/golang-jwt/jwt/v5"

// JWT
type contextKey string

const UserIDkey contextKey = "userID"

type Claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}

// Auth Request
type AuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthDTO struct {
	Token string `json:"token"`
}
