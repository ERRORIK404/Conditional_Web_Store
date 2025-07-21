package config

import (
	"fmt"

	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

type config struct {
	PostgresHost string `env:"POSTGRES_HOST" envDefault:"db"`
	PostgresPort string `env:"POSTGRES_PORT" envDefault:"5432"`

	PostgresDB string `env:"POSTGRES_DB" envDefault:"marketplace"`
	PostgresUser string `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	

	ApiPort string `env:"API_PORT"`

	DB_URL string
}

var Сfg config

func LoadConfig() (*config, error) {
	err := env.Parse(&Сfg)
	if err != nil {
		logger.Logger.Error("Error loading config", zap.Error(err))
		return nil, err
	}

	Сfg.DB_URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", Сfg.PostgresUser, Сfg.PostgresPassword, Сfg.PostgresHost, Сfg.PostgresPort, Сfg.PostgresDB)
	logger.Logger.Debug("Config loaded", zap.Any("config", Сfg))

	return &Сfg, nil
}

