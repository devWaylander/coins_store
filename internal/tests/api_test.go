package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/devWaylander/coins_store/internal/handler"
	"github.com/devWaylander/coins_store/internal/middleware/auth"
	"github.com/devWaylander/coins_store/internal/middleware/cors"
	"github.com/devWaylander/coins_store/internal/middleware/logger"
	"github.com/devWaylander/coins_store/internal/repo"
	"github.com/devWaylander/coins_store/internal/service"
	"github.com/devWaylander/coins_store/pkg/log"
	"github.com/devWaylander/coins_store/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/devWaylander/coins_store/config"
)

type E2eIntegrationTestSuite struct {
	suite.Suite
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
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Logger.Fatal().Msgf("Database connection failed: %v", err)
	}

	log.Logger.Info().Msg("Database connection established successfully")

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

	httpServer.ListenAndServe()
}

func (s *E2eIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()

	cfg, err := config.Parse()
	if err != nil {
		log.Logger.Fatal().Msg(err.Error())
	}

	dbConfig, err := pgxpool.ParseConfig(cfg.DB.DBTestUrl)
	if err != nil {
		log.Logger.Fatal().Msgf("Unable to parse database URL: %v\n", err)
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Logger.Fatal().Msgf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Logger.Fatal().Msgf("Database connection failed: %v", err)
	}

	tablesToClear := []string{
		"shop.balance",
		"shop.balance_history",
		"shop.user",
		"shop.inventory",
		"shop.inventory_merch",
	}

	for _, table := range tablesToClear {
		_, err := dbPool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table))
		if err != nil {
			log.Logger.Fatal().Msgf("Failed to truncate table %s: %v\n", table, err)
		}
	}

	log.Logger.Info().Msg("Database cleared successfully (except for merch table).")
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
	t.Run("fail_buy_third_merch", func(t *testing.T) {
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
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

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
