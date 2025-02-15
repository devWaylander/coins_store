package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	internalErrors "github.com/devWaylander/coins_store/pkg/errors"
	"github.com/devWaylander/coins_store/pkg/log"
	"github.com/devWaylander/coins_store/pkg/models"
)

type AuthMiddleware interface {
	LoginWithPass(ctx context.Context, qp models.AuthQuery) (models.AuthDTO, error)
}

type Service interface {
	GetUserInfo(ctx context.Context, qp models.InfoQuery) (models.InfoDTO, error)
	BuyItem(ctx context.Context, qp models.ItemQuery) error
	SendCoins(ctx context.Context, qp models.CoinsQuery) error
}

func New(ctx context.Context, mux *http.ServeMux, authMiddleware AuthMiddleware, service Service) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request: Method not allowed", http.StatusBadRequest)
	})

	// unsecured handles
	// Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
	mux.HandleFunc("POST /api/auth", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body := models.AuthReqBody{}

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

		authDTO, err := authMiddleware.LoginWithPass(ctx, models.AuthQuery(body))
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
		ctx := r.Context()

		claims, err := decodeCtxClaims(ctx)
		if err != nil {
			http.Error(w, internalErrors.ErrGetInfo, http.StatusInternalServerError)
		}

		infoDTO, err := service.GetUserInfo(ctx, models.InfoQuery{
			UserID:   claims.UserID,
			Username: claims.Username,
		})
		if err != nil {
			log.Logger.Err(err).Msg(err.Error())
			http.Error(w, internalErrors.ErrGetInfo, http.StatusInternalServerError)
			return
		}

		sendResponse(w, infoDTO)
	})
	// Купить предмет за монеты.
	mux.HandleFunc("GET /api/buy/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		item := r.URL.Path[len("/api/buy/"):]
		if item == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		claims, err := decodeCtxClaims(ctx)
		if err != nil {
			http.Error(w, internalErrors.ErrGetInfo, http.StatusInternalServerError)
		}

		err = service.BuyItem(ctx, models.ItemQuery{
			UserID:   claims.UserID,
			Username: claims.Username,
			Item:     item,
		})
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
	// Отправить монеты другому пользователю
	mux.HandleFunc("POST /api/sendCoin", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body := models.SendCoinsReqBody{}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Logger.Err(err).Msg(err.Error())
			http.Error(w, internalErrors.ErrUnmarshalResponse, http.StatusInternalServerError)
			return
		}
		if body.Amount < 1 {
			http.Error(w, internalErrors.ErrInvalidSendCoinsReqParams, http.StatusBadRequest)
			return
		}

		claims, err := decodeCtxClaims(ctx)
		if err != nil {
			http.Error(w, internalErrors.ErrGetInfo, http.StatusInternalServerError)
		}
		if body.Recipient == claims.Username {
			http.Error(w, internalErrors.ErrInvalidRecipientYourself, http.StatusBadRequest)
			return
		}

		err = service.SendCoins(ctx, models.CoinsQuery{
			UserID:    claims.UserID,
			Amount:    body.Amount,
			Sender:    claims.Username,
			Recipient: body.Recipient,
		})
		if err != nil {
			switch err.Error() {
			case internalErrors.ErrInvalidRecipient:
				http.Error(w, internalErrors.ErrInvalidRecipient, http.StatusBadRequest)
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

func decodeCtxClaims(ctx context.Context) (models.Claims, error) {
	err := errors.New(internalErrors.ErrDecodeCtx)

	userID, ok := ctx.Value(models.UserIDKey).(int64)
	if !ok {
		log.Logger.Err(err).Msg(err.Error())
		return models.Claims{}, err
	}
	username, ok := ctx.Value(models.UsernameKey).(string)
	if !ok {
		log.Logger.Err(err).Msg(err.Error())
		return models.Claims{}, err
	}

	return models.Claims{UserID: userID, Username: username}, nil
}
