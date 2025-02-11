package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/devWaylander/coins_store/pkg/log"
	"github.com/joho/godotenv"
)

var (
	C Config
)

type Config struct {
	Common Common `envPrefix:"COMMON_"`
	DB     DB     `envPrefix:"DB_"`
}

type Common struct {
	Port             string `env:"PORT,required"`
	JWTSecret        string `env:"JWT_SECRET,required"`
	JWTRefreshSecret string `env:"JWT_REFRESH_SECRET,required"`
}

type DB struct {
	DBHost           string `env:"DB_HOST,required"`
	DBDriver         string `env:"DB_DRIVER,required"`
	DBUser           string `env:"DB_USER,required"`
	DBPassword       string `env:"DB_PASSWORD,required"`
	DBName           string `env:"DB_NAME,required"`
	DBPort           string `env:"DB_PORT,required"`
	DBMaxConnections int32  `env:"DB_MAX_CONNECTIONS,required"`
	DBUrl            string `env:"DATABASE_URL,required"`
	TestDBUrl        string `env:"TEST_DATABASE_URL,required"`
}

func Parse() (Config, error) {
	// Read envs
	err := godotenv.Load("../../.env.local")
	if err != nil {
		err = godotenv.Load("../../.env")
		if err != nil {
			return C, fmt.Errorf("failed to read environment variables: %w", err)
		}
	}

	// Decode envs into config structures
	err = env.Parse(&C)
	if err != nil {
		if aggErr, ok := err.(env.AggregateError); ok {
			for _, e := range aggErr.Errors {
				log.Logger.Error().Msg(fmt.Sprintf("Validation error: '%s'\n", e.Error()))
			}
		}

		return C, err
	}

	return C, nil
}
