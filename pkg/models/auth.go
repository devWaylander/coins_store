package models

import "github.com/golang-jwt/jwt/v5"

// JWT
type contextKey string

const UserIDKey contextKey = "userID"
const UsernameKey contextKey = "username"

type Claims struct {
	UserID   int64  `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthDTO struct {
	Token string `json:"token"`
}
