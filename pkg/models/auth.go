package models

import "github.com/golang-jwt/jwt/v5"

// JWT
type contextKey string

const UsernameKey contextKey = "username"

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Auth Request
type AuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
