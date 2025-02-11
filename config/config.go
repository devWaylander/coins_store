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
	Port      string `env:"API_PORT,required"`
	JWTSecret string `env:"JWT_SECRET,required"`
}

type DB struct {
	DBHost           string `env:"HOST,required"`
	DBUser           string `env:"USER,required"`
	DBPassword       string `env:"PASSWORD,required"`
	DBName           string `env:"NAME,required"`
	DBPort           string `env:"PORT,required"`
	DBMaxConnections int32  `env:"MAX_CONNECTIONS,required"`
	DBUrl            string `env:"DATABASE_URL,required"`
	TestDBUrl        string `env:"TEST_DATABASE_URL,required"`
}

func Parse() (Config, error) {
	// Read envs
	err := godotenv.Load("../.env.local")
	if err != nil {
		err = godotenv.Load("../.env")
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
