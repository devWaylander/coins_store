package middleware

import (
	"context"
	"net/http"

	"github.com/devWaylander/coins_store/pkg/errors"
	"github.com/devWaylander/coins_store/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

var unsecuredHandles = map[string]struct{}{
	"/api/auth": {},
}

func AuthMiddleware(next http.Handler, jwtKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if unsecuredHandles[r.URL.Path] == struct{}{} {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, errors.ErrAuthHeader, http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, errors.ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			http.Error(w, errors.ErrInvalidClaims, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), models.UserIDkey, claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
