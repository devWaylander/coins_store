package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/devWaylander/coins_store/config"
	"github.com/devWaylander/coins_store/internal/handler"
	auth "github.com/devWaylander/coins_store/internal/middleware/auth"
	logger "github.com/devWaylander/coins_store/internal/middleware/logger"
	"github.com/devWaylander/coins_store/internal/repo"
	"github.com/devWaylander/coins_store/internal/service"
	errorgroup "github.com/devWaylander/coins_store/pkg/error_group"
	"github.com/devWaylander/coins_store/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Config
	cfg, err := config.Parse()
	if err != nil {
		log.Logger.Fatal().Msg(err.Error())
	}

	// Graceful shutdown init
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	// DB
	dbConfig, err := pgxpool.ParseConfig(cfg.DB.DBUrl)
	if err != nil {
		log.Logger.Fatal().Msgf("Unable to parse database URL: %v\n", err)
	}

	dbConfig.MaxConns = cfg.DB.DBMaxConnections
	dbConfig.MaxConnLifetime = 1 * time.Minute
	dbConfig.MaxConnIdleTime = 1 * time.Minute

	dbPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Logger.Fatal().Msgf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	// Migrations
	cmd := exec.Command("dbmate", "-u", cfg.DB.DBUrl, "--migrations-dir", "../db/migrations", "--no-dump-schema", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Logger.Fatal().Msgf("Error running dbmate: %v", err)
	}

	// Repositories
	usecaseRepo := repo.New(dbPool)
	authMiddlewareRepo := auth.NewAuthRepo(dbPool)

	// Auth Middleware
	authMiddleware := auth.NewMiddleware(authMiddlewareRepo, cfg.Common.JWTSecret)

	// Service
	service := service.New(usecaseRepo, cfg.Common.JWTSecret)

	// Handler
	mux := http.NewServeMux()
	handler.New(ctx, mux, authMiddleware, service)
	wrappedAuthMux := authMiddleware.Middleware(mux)
	wrappedLoggerMux := logger.Middleware(wrappedAuthMux)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Common.Port),
		Handler: wrappedLoggerMux,
	}

	// Graceful shutdown run
	g, gCtx := errorgroup.EGWithContext(ctx)
	g.Go(func() error {
		log.Logger.Info().Msgf("Server is up on port: %s", cfg.Common.Port)
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Logger.Info().Msgf("Server on port %s is shutting down", cfg.Common.Port)
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		log.Logger.Info().Msgf("exit reason: %s \\n", err)
	}
}
