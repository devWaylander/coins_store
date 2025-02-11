package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type contextKey string

const UsernameKey contextKey = "username"
