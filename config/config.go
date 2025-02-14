package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/devWaylander/coins_store/pkg/log"
	"github.com/joho/godotenv"
)

var (
	C           Config
	isContainer bool = true
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
	DBUrl            string `json:"databaseURL"`
	DBLocalUrl       string `env:"DATABASE_LOCAL_URL,required"`
	DBContainerUrl   string `env:"DATABASE_CONTAINER_URL,required"`
	TestDBUrl        string `env:"TEST_DATABASE_URL,required"`
}

func Parse() (Config, error) {
	// If running from container, use docker envs
	if _, exists := os.LookupEnv("COMMON_API_PORT"); !exists {
		// If running on local machine use env file
		err := godotenv.Load("../.env")
		if err != nil {
			return C, fmt.Errorf("failed to read environment variables: %w", err)
		}
		isContainer = false
	}

	// Decode envs into config structures
	err := env.Parse(&C)
	if err != nil {
		if aggErr, ok := err.(env.AggregateError); ok {
			for _, e := range aggErr.Errors {
				log.Logger.Error().Msg(fmt.Sprintf("Validation error: '%s'\n", e.Error()))
			}
		}

		return C, err
	}

	C.DB.DBUrl = C.DB.DBLocalUrl
	if isContainer {
		C.DB.DBUrl = C.DB.DBContainerUrl
	}

	return C, nil
}
