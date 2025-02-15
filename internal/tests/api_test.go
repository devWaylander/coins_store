package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/devWaylander/coins_store/internal/handler"
	"github.com/devWaylander/coins_store/internal/middleware/auth"
	"github.com/devWaylander/coins_store/internal/middleware/cors"
	"github.com/devWaylander/coins_store/internal/middleware/logger"
	"github.com/devWaylander/coins_store/internal/repo"
	"github.com/devWaylander/coins_store/internal/service"
	internalErrors "github.com/devWaylander/coins_store/pkg/errors"
	"github.com/devWaylander/coins_store/pkg/log"
	"github.com/devWaylander/coins_store/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/devWaylander/coins_store/config"
)

type E2eIntegrationTestSuite struct {
	suite.Suite
	dbPool *pgxpool.Pool
}

func TestE2eIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(E2eIntegrationTestSuite))
}

func (s *E2eIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()

	// Config
	cfg, err := config.Parse()
	if err != nil {
		log.Logger.Fatal().Msg(err.Error())
	}

	dbConfig, err := pgxpool.ParseConfig(cfg.DB.DBTestUrl)
	if err != nil {
		log.Logger.Fatal().Msgf("Unable to parse database URL: %v\n", err)
	}

	dbConfig.MaxConns = cfg.DB.DBMaxConnections
	dbConfig.MaxConnLifetime = cfg.DB.DBLifeTimeConnection
	dbConfig.MaxConnIdleTime = cfg.DB.DBMaxConnIdleTime

	dbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Logger.Fatal().Msgf("Unable to create connection pool: %v\n", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		log.Logger.Fatal().Msgf("Database connection failed: %v", err)
	}

	log.Logger.Info().Msg("Database connection established successfully")
	s.dbPool = dbPool

	// Repositories
	usecaseRepo := repo.New(dbPool)
	authMiddlewareRepo := auth.NewAuthRepo(dbPool)

	// Auth Middleware
	authMiddleware := auth.NewMiddleware(authMiddlewareRepo, cfg.Common.JWTSecret)

	// Service
	service := service.New(usecaseRepo)

	// Handler
	mux := http.NewServeMux()
	handler.New(ctx, mux, authMiddleware, service)
	wrappedAuthMux := authMiddleware.Middleware(mux)
	wrappedCorsMux := cors.Middleware(wrappedAuthMux)
	wrappedLoggerMux := logger.Middleware(wrappedCorsMux)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Common.Port),
		Handler: wrappedLoggerMux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Fatal().Msgf("HTTP server error: %v", err)
		}
	}()
}

func (s *E2eIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()

	tablesToClear := []string{
		"shop.balance",
		"shop.balance_history",
		"shop.user",
		"shop.inventory",
		"shop.inventory_merch",
	}

	for _, table := range tablesToClear {
		_, err := s.dbPool.Exec(ctx, fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE;`, table))
		if err != nil {
			log.Logger.Fatal().Msgf("Failed to truncate table %s: %v\n", table, err)
		}
	}

	log.Logger.Info().Msg("Database cleared successfully")
	s.dbPool.Close()
}

func (s *E2eIntegrationTestSuite) TestAuth() {
	t := s.T()
	client := HttpClient{}

	t.Run("success_new_user_auth", func(t *testing.T) {
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user1",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respData := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respData)
		require.NoError(t, err)

		require.Greater(t, len(respData.Token), 0)
	})

	t.Run("success_exist_user_auth", func(t *testing.T) {
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user1",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respData := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respData)
		require.NoError(t, err)

		require.Greater(t, len(respData.Token), 0)
	})

	t.Run("fail_exist_user_auth", func(t *testing.T) {
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user1",
			Password: "1asdfjnaasf",
		})
		require.NoError(t, err)

		resp, _, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func (s *E2eIntegrationTestSuite) TestBuyMerch() {
	t := s.T()
	client := HttpClient{}

	t.Run("success_buy_merch", func(t *testing.T) {
		// login
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user1",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthData := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthData)
		require.NoError(t, err)

		require.Greater(t, len(respAuthData.Token), 0)

		// buy
		resp, _, err = client.SendJsonReq(respAuthData.Token, http.MethodGet, BaseURL+"/api/buy/pink-hoody", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// get info with item
		resp, respBody, err = client.SendJsonReq(respAuthData.Token, http.MethodGet, BaseURL+"/api/info", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respInfoData := models.InfoDTO{}
		err = json.Unmarshal(respBody, &respInfoData)
		require.NoError(t, err)

		expectedInfoData := models.InfoDTO{Inventory: []models.MerchDTO{{Type: "pink-hoody", Quantity: 1}}}
		for _, item := range respInfoData.Inventory {
			require.Equal(t, expectedInfoData.Inventory[0].Type, item.Type)
			require.Equal(t, expectedInfoData.Inventory[0].Quantity, item.Quantity)
		}
	})
	t.Run("success_buy_second_merch", func(t *testing.T) {
		// login
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user1",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthData := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthData)
		require.NoError(t, err)

		require.Greater(t, len(respAuthData.Token), 0)

		// buy
		resp, _, err = client.SendJsonReq(respAuthData.Token, http.MethodGet, BaseURL+"/api/buy/pink-hoody", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// get info with item
		resp, respBody, err = client.SendJsonReq(respAuthData.Token, http.MethodGet, BaseURL+"/api/info", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respInfoData := models.InfoDTO{}
		err = json.Unmarshal(respBody, &respInfoData)
		require.NoError(t, err)

		expectedInfoData := models.InfoDTO{Inventory: []models.MerchDTO{{Type: "pink-hoody", Quantity: 2}}}
		for _, item := range respInfoData.Inventory {
			require.Equal(t, expectedInfoData.Inventory[0].Type, item.Type)
			require.Equal(t, expectedInfoData.Inventory[0].Quantity, item.Quantity)
		}
	})
	t.Run("fail_buy_third_merch_not_enough_coins", func(t *testing.T) {
		// login
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user1",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthData := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthData)
		require.NoError(t, err)

		require.Greater(t, len(respAuthData.Token), 0)

		// buy
		resp, respBody, err = client.SendJsonReq(respAuthData.Token, http.MethodGet, BaseURL+"/api/buy/pink-hoody", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, internalErrors.ErrNotEnoughCoins, strings.TrimSpace(string(respBody)))

		// get info with item
		resp, respBody, err = client.SendJsonReq(respAuthData.Token, http.MethodGet, BaseURL+"/api/info", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respInfoData := models.InfoDTO{}
		err = json.Unmarshal(respBody, &respInfoData)
		require.NoError(t, err)

		expectedInfoData := models.InfoDTO{Inventory: []models.MerchDTO{{Type: "pink-hoody", Quantity: 2}}}
		for _, item := range respInfoData.Inventory {
			require.Equal(t, expectedInfoData.Inventory[0].Type, item.Type)
			require.Equal(t, expectedInfoData.Inventory[0].Quantity, item.Quantity)
		}
	})
	t.Run("fail_buy_third_merch_doesn't_exist", func(t *testing.T) {
		// login
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user1",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthData := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthData)
		require.NoError(t, err)

		require.Greater(t, len(respAuthData.Token), 0)

		// buy
		resp, respBody, err = client.SendJsonReq(respAuthData.Token, http.MethodGet, BaseURL+"/api/buy/purple-hoody", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, internalErrors.ErrItemDoesntExist, strings.TrimSpace(string(respBody)))

		// get info with item
		resp, respBody, err = client.SendJsonReq(respAuthData.Token, http.MethodGet, BaseURL+"/api/info", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respInfoData := models.InfoDTO{}
		err = json.Unmarshal(respBody, &respInfoData)
		require.NoError(t, err)

		expectedInfoData := models.InfoDTO{Inventory: []models.MerchDTO{{Type: "pink-hoody", Quantity: 2}}}
		for _, item := range respInfoData.Inventory {
			require.Equal(t, expectedInfoData.Inventory[0].Type, item.Type)
			require.Equal(t, expectedInfoData.Inventory[0].Quantity, item.Quantity)
		}
	})
}

func (s *E2eIntegrationTestSuite) TestSendCoins() {
	t := s.T()
	client := HttpClient{}

	t.Run("success_send_coins", func(t *testing.T) {
		// create 1 user
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user2",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthDataU1 := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthDataU1)
		require.NoError(t, err)

		require.Greater(t, len(respAuthDataU1.Token), 0)

		// create 2 user
		reqBody, err = json.Marshal(models.AuthReqBody{
			Username: "user3",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err = client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthDataU2 := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthDataU2)
		require.NoError(t, err)

		require.Greater(t, len(respAuthDataU2.Token), 0)

		// send coins from user 1 to user 2
		reqBody, err = json.Marshal(models.SendCoinsReqBody{
			Recipient: "user3",
			Amount:    69,
		})
		require.NoError(t, err)
		resp, _, err = client.SendJsonReq(respAuthDataU1.Token, http.MethodPost, BaseURL+"/api/sendCoin", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// user1 info
		resp, respBody, err = client.SendJsonReq(respAuthDataU1.Token, http.MethodGet, BaseURL+"/api/info", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respInfoData := models.InfoDTO{}
		err = json.Unmarshal(respBody, &respInfoData)
		require.NoError(t, err)

		expectedInfoData := models.InfoDTO{
			Coins: 931,
			CoinsHistory: models.BalanceHistoryDTO{
				Received: []models.ReceivedDTO{},
				Sent:     []models.SentDTO{{ToUser: "user3", Amount: 69}},
			},
		}
		require.Equal(t, expectedInfoData.Coins, respInfoData.Coins)
		require.ElementsMatch(t, expectedInfoData.CoinsHistory.Sent, respInfoData.CoinsHistory.Sent)
		require.ElementsMatch(t, expectedInfoData.CoinsHistory.Received, respInfoData.CoinsHistory.Received)

		// user2 info
		resp, respBody, err = client.SendJsonReq(respAuthDataU2.Token, http.MethodGet, BaseURL+"/api/info", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respInfoData = models.InfoDTO{}
		err = json.Unmarshal(respBody, &respInfoData)
		require.NoError(t, err)

		expectedInfoData = models.InfoDTO{
			Coins: 1069,
			CoinsHistory: models.BalanceHistoryDTO{
				Received: []models.ReceivedDTO{{FromUser: "user2", Amount: 69}},
				Sent:     []models.SentDTO{},
			},
		}
		require.Equal(t, expectedInfoData.Coins, respInfoData.Coins)
		require.ElementsMatch(t, expectedInfoData.CoinsHistory.Sent, respInfoData.CoinsHistory.Sent)
		require.ElementsMatch(t, expectedInfoData.CoinsHistory.Received, respInfoData.CoinsHistory.Received)
	})
	t.Run("fail_send_yourself_coins", func(t *testing.T) {
		// login 1 user
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user2",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthDataU1 := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthDataU1)
		require.NoError(t, err)

		require.Greater(t, len(respAuthDataU1.Token), 0)

		// send coins from user 1 to user 2
		reqBody, err = json.Marshal(models.SendCoinsReqBody{
			Recipient: "user2",
			Amount:    100,
		})
		require.NoError(t, err)
		resp, respBody, err = client.SendJsonReq(respAuthDataU1.Token, http.MethodPost, BaseURL+"/api/sendCoin", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, internalErrors.ErrInvalidRecipientYourself, strings.TrimSpace(string(respBody)))
	})
	t.Run("fail_send_doesn't_exist_user_coins", func(t *testing.T) {
		// login 1 user
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user2",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthDataU1 := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthDataU1)
		require.NoError(t, err)

		require.Greater(t, len(respAuthDataU1.Token), 0)

		// send coins from user 1 to user 2
		reqBody, err = json.Marshal(models.SendCoinsReqBody{
			Recipient: "user5468974984",
			Amount:    100,
		})
		require.NoError(t, err)
		resp, respBody, err = client.SendJsonReq(respAuthDataU1.Token, http.MethodPost, BaseURL+"/api/sendCoin", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, internalErrors.ErrInvalidRecipient, strings.TrimSpace(string(respBody)))
	})
	t.Run("success_second_send_coins", func(t *testing.T) {
		// login 1 user
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user2",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthDataU1 := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthDataU1)
		require.NoError(t, err)

		require.Greater(t, len(respAuthDataU1.Token), 0)

		// login 2 user
		reqBody, err = json.Marshal(models.AuthReqBody{
			Username: "user3",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err = client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthDataU2 := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthDataU2)
		require.NoError(t, err)

		require.Greater(t, len(respAuthDataU2.Token), 0)

		// send coins from user 1 to user 2
		reqBody, err = json.Marshal(models.SendCoinsReqBody{
			Recipient: "user3",
			Amount:    931,
		})
		require.NoError(t, err)
		resp, _, err = client.SendJsonReq(respAuthDataU1.Token, http.MethodPost, BaseURL+"/api/sendCoin", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// user1 info
		resp, respBody, err = client.SendJsonReq(respAuthDataU1.Token, http.MethodGet, BaseURL+"/api/info", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respInfoData := models.InfoDTO{}
		err = json.Unmarshal(respBody, &respInfoData)
		require.NoError(t, err)

		expectedInfoData := models.InfoDTO{
			Coins: 0,
			CoinsHistory: models.BalanceHistoryDTO{
				Received: []models.ReceivedDTO{},
				Sent:     []models.SentDTO{{ToUser: "user3", Amount: 69}, {ToUser: "user3", Amount: 931}},
			},
		}
		require.Equal(t, expectedInfoData.Coins, respInfoData.Coins)
		require.ElementsMatch(t, expectedInfoData.CoinsHistory.Sent, respInfoData.CoinsHistory.Sent)
		require.ElementsMatch(t, expectedInfoData.CoinsHistory.Received, respInfoData.CoinsHistory.Received)

		// user2 info
		resp, respBody, err = client.SendJsonReq(respAuthDataU2.Token, http.MethodGet, BaseURL+"/api/info", []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respInfoData = models.InfoDTO{}
		err = json.Unmarshal(respBody, &respInfoData)
		require.NoError(t, err)

		expectedInfoData = models.InfoDTO{
			Coins: 2000,
			CoinsHistory: models.BalanceHistoryDTO{
				Received: []models.ReceivedDTO{{FromUser: "user2", Amount: 69}, {FromUser: "user2", Amount: 931}},
				Sent:     []models.SentDTO{},
			},
		}
		require.Equal(t, expectedInfoData.Coins, respInfoData.Coins)
		require.ElementsMatch(t, expectedInfoData.CoinsHistory.Sent, respInfoData.CoinsHistory.Sent)
		require.ElementsMatch(t, expectedInfoData.CoinsHistory.Received, respInfoData.CoinsHistory.Received)
	})
	t.Run("fail_third_send_coins", func(t *testing.T) {
		// login 1 user
		reqBody, err := json.Marshal(models.AuthReqBody{
			Username: "user2",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err := client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthDataU1 := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthDataU1)
		require.NoError(t, err)

		require.Greater(t, len(respAuthDataU1.Token), 0)

		// login 2 user
		reqBody, err = json.Marshal(models.AuthReqBody{
			Username: "user3",
			Password: "11111!Aa",
		})
		require.NoError(t, err)

		resp, respBody, err = client.SendJsonReq("", http.MethodPost, BaseURL+"/api/auth", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respAuthDataU2 := models.AuthDTO{}
		err = json.Unmarshal(respBody, &respAuthDataU2)
		require.NoError(t, err)

		require.Greater(t, len(respAuthDataU2.Token), 0)

		// send coins from user 1 to user 2
		reqBody, err = json.Marshal(models.SendCoinsReqBody{
			Recipient: "user3",
			Amount:    100,
		})
		require.NoError(t, err)
		resp, respBody, err = client.SendJsonReq(respAuthDataU1.Token, http.MethodPost, BaseURL+"/api/sendCoin", reqBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, internalErrors.ErrNotEnoughCoins, strings.TrimSpace(string(respBody)))
	})
}

// func (s *E2eIntegrationTestSuite) TestGetUserInfo() {
// 	t := s.T()
// 	client := HttpClient{}

// 	t.Run("", func(t *testing.T) {

// 	})
// }
