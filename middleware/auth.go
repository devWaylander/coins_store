package middleware

import (
	"context"
	"net/http"

	"github.com/devWaylander/coins_store/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler, jwtKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			http.Error(w, "Couldn't parse claims", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), models.UsernameKey, claims.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
