package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	internalErrors "github.com/devWaylander/coins_store/pkg/errors"
	"github.com/devWaylander/coins_store/pkg/log"
	"github.com/devWaylander/coins_store/pkg/models"
)

type AuthMiddleware interface {
	LoginWithPass(ctx context.Context, username, password string) (models.AuthDTO, error)
}

type Service interface {
}

func New(ctx context.Context, mux *http.ServeMux, authMiddleware AuthMiddleware, service Service) {
	// unsecured handles
	mux.HandleFunc("POST /api/auth", func(w http.ResponseWriter, r *http.Request) {
		body := models.AuthReq{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Logger.Err(err)
			http.Error(w, internalErrors.ErrUnmarshalResponse, http.StatusInternalServerError)
			return
		}
		if body.Password == "" || body.Username == "" {
			log.Logger.Info().Msg(fmt.Sprintf("Empty request: %s", internalErrors.ErrInvalidAuthReqParams))
			http.Error(w, internalErrors.ErrInvalidAuthReqParams, http.StatusBadRequest)
			return
		}

		authDTO, err := authMiddleware.LoginWithPass(ctx, body.Username, body.Password)
		if err != nil {
			if internalErrors.ErrWrongPassword == err.Error() {
				http.Error(w, internalErrors.ErrWrongPassword, http.StatusUnauthorized)
				return
			}

			log.Logger.Err(err)
			http.Error(w, internalErrors.ErrLogin, http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(authDTO)
		if err != nil {
			log.Logger.Err(err)
			http.Error(w, internalErrors.ErrMarshalResponse, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(data); err != nil {
			log.Logger.Err(err)
			http.Error(w, internalErrors.ErrMarshalResponse, http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request: Method not allowed", http.StatusBadRequest)
	})

	// secured handles
}
