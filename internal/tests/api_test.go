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
}

func (s *E2eIntegrationTestSuite) TestAuth() {
	t := s.T()
	client := HttpClient{}

	t.Run("new_user_auth_success", func(t *testing.T) {
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

	t.Run("exist_user_auth_success", func(t *testing.T) {
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
}
