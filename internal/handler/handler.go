package handler

import (
	"context"
	"net/http"
)

type Service interface {
	GenerateJWT(username string) (string, error)
}

func New(ctx context.Context, mux *http.ServeMux, service Service) {
	// unsecured handles
	mux.HandleFunc("POST /api/auth", func(w http.ResponseWriter, r *http.Request) {
		// распарсить body из запроса username, password
		// пойти проверить не существует ли такого username
		// если существует, то проверить его пароль
		// если не существует, то сгенерить токен, создать пользователя с паролем, монетами, инвентарём
		// service.GenerateJWT()
	})

	// secured handles
}
