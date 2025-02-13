package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	internalErrors "github.com/devWaylander/coins_store/pkg/errors"
	"github.com/devWaylander/coins_store/pkg/log"
	"github.com/devWaylander/coins_store/pkg/models"
)

type AuthMiddleware interface {
	LoginWithPass(ctx context.Context, username, password string) (models.AuthDTO, error)
}

type Service interface {
	GetUserInfo(ctx context.Context, userID int64, username string) (models.InfoDTO, error)
	BuyItem(ctx context.Context, userID int64, username, item string) error
}

func New(ctx context.Context, mux *http.ServeMux, authMiddleware AuthMiddleware, service Service) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request: Method not allowed", http.StatusBadRequest)
	})

	// unsecured handles
	// Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
	mux.HandleFunc("POST /api/auth", func(w http.ResponseWriter, r *http.Request) {
		body := models.AuthReq{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Logger.Err(err).Msg(err.Error())
			http.Error(w, internalErrors.ErrUnmarshalResponse, http.StatusInternalServerError)
			return
		}
		if body.Password == "" || body.Username == "" {
			http.Error(w, internalErrors.ErrInvalidAuthReqParams, http.StatusBadRequest)
			return
		}

		authDTO, err := authMiddleware.LoginWithPass(ctx, body.Username, body.Password)
		if err != nil {
			switch err.Error() {
			case internalErrors.ErrWrongPassword:
				http.Error(w, internalErrors.ErrWrongPassword, http.StatusUnauthorized)
			case internalErrors.ErrWrongPasswordFormat:
				http.Error(w, internalErrors.ErrWrongPasswordFormat, http.StatusUnauthorized)
			default:
				http.Error(w, internalErrors.ErrLogin, http.StatusInternalServerError)
				log.Logger.Err(err).Msg(err.Error())
			}
			return
		}

		sendResponse(w, authDTO)
	})

	// secured handles
	// Получить информацию о монетах, инвентаре и истории транзакций.
	mux.HandleFunc("GET /api/info", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(models.UserIDKey).(int64)
		username := r.Context().Value(models.UsernameKey).(string)
		infoDTO, err := service.GetUserInfo(ctx, userID, username)
		if err != nil {
			log.Logger.Err(err).Msg(err.Error())
			http.Error(w, internalErrors.ErrGetInfo, http.StatusInternalServerError)
			return
		}

		sendResponse(w, infoDTO)
	})
	// Купить предмет за монеты.
	mux.HandleFunc("GET /api/buy/", func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")
		if len(urlParts) < 4 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		item := urlParts[3]
		if item == "" {
			http.Error(w, internalErrors.ErrInvalidGetBuyItemReqParams, http.StatusBadRequest)
			return
		}

		userID := r.Context().Value(models.UserIDKey).(int64)
		username := r.Context().Value(models.UsernameKey).(string)
		err := service.BuyItem(ctx, userID, username, item)
		if err != nil {
			switch err.Error() {
			case internalErrors.ErrItemDoesntExist:
				http.Error(w, internalErrors.ErrItemDoesntExist, http.StatusBadRequest)
			case internalErrors.ErrNotEnoughCoins:
				http.Error(w, internalErrors.ErrNotEnoughCoins, http.StatusBadRequest)
			default:
				http.Error(w, internalErrors.ErrGetBuyItem, http.StatusInternalServerError)
				log.Logger.Err(err).Msg(err.Error())
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func sendResponse(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Logger.Err(err).Msg(err.Error())
		http.Error(w, internalErrors.ErrMarshalResponse, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		log.Logger.Err(err).Msg(err.Error())
		http.Error(w, internalErrors.ErrMarshalResponse, http.StatusInternalServerError)
		return
	}
}
