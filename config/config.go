package config

import (
	"fmt"
	"os"
	"time"

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
	DBHost               string        `env:"HOST,required"`
	DBUser               string        `env:"USER,required"`
	DBPassword           string        `env:"PASSWORD,required"`
	DBName               string        `env:"NAME,required"`
	DBPort               string        `env:"PORT,required"`
	DBMaxConnections     int32         `env:"MAX_CONNECTIONS,required"`
	DBLocalUrl           string        `env:"DATABASE_LOCAL_URL,required"`
	DBContainerUrl       string        `env:"DATABASE_CONTAINER_URL,required"`
	DBLifeTimeConnection time.Duration `json:"dbLifeTimeConnection"`
	DBMaxConnIdleTime    time.Duration `json:"dbeMaxIdleTime"`
	DBUrl                string        `json:"dbURL"`
}

func Parse() (Config, error) {
	isContainer := isRunningInContainer()

	if !isContainer {
		err := godotenv.Load("../.env")
		if err != nil {
			return C, fmt.Errorf("failed to read environment variables: %w", err)
		}
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

	C.DB.DBLifeTimeConnection = 1 * time.Minute
	C.DB.DBMaxConnIdleTime = 1 * time.Minute

	return C, nil
}

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}
