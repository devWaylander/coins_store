package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devWaylander/coins_store/pkg/errors"
	"github.com/devWaylander/coins_store/pkg/log"
	"github.com/devWaylander/coins_store/pkg/models"
)

type Service interface {
	ValidateUser(ctx context.Context, username, password string) error
}

func New(ctx context.Context, mux *http.ServeMux, service Service) {
	// unsecured handles
	mux.HandleFunc("POST /api/auth", func(w http.ResponseWriter, r *http.Request) {
		// пойти проверить не существует ли такого username
		// если существует, то проверить его пароль
		// если пароль неверный вернуть 401
		// если не существует, то сгенерить токен, создать пользователя с паролем, монетами, инвентарём
		// в обоих случаях успеха отдать токен
		body := models.AuthReq{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Logger.Err(err)
			http.Error(w, errors.ErrUnmarshalResponse, http.StatusInternalServerError)
		}
		if body.Password == "" || body.Username == "" {
			log.Logger.Info().Msg(fmt.Sprintf("Empty request: %s", errors.ErrInvalidAuthReqParams))
			http.Error(w, errors.ErrInvalidAuthReqParams, http.StatusBadRequest)
		}

		// service.GenerateJWT()
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request: Method not allowed", http.StatusBadRequest)
	})

	// secured handles
}
