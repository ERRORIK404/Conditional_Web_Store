package main

import (
	"time"

	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/config"
	localjwt "github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/localJWT"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger(false)
	_, err := config.LoadConfig()
	if err != nil {
		logger.Logger.Error("Error loading config", zap.Error(err))
		return
	}

	var jwtSecret struct {
		secret string `env:"JWT_SECRET"`
	}
	err = env.Parse(&jwtSecret)
	if err != nil {
		logger.Logger.Error("Error loading JWT secret", zap.Error(err))
		return
	}
	localjwt.LoadJWT(jwtSecret.secret)

	time.Sleep(150 * time.Second)
}
