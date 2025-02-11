package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devWaylander/coins_store/config"
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

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Logger.Fatal().Msgf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	// Repository
	// Service
	// Handler

	// Graceful shutdown run
	// g, gCtx := errgroup.Group
	// g.Go(func() error {
	// 	log.Logger.Info().Msgf("Server is up on: %s:%s", cfg.Common., port)
	// 	return httpServer.ListenAndServe()
	// })
	// g.Go(func() error {
	// 	<-gCtx.Done()
	// 	log.Logger.Info().Msgf("Server is shutting down: %s:%s", ip, port)
	// 	return httpServer.Shutdown(context.Background())
	// })

	// if _, err := g.Wait(); err != nil {
	// 	log.Logger.Info().Msgf("exit reason: %s \\n", err)
	// }
}
