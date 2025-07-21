package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
    Logger *zap.Logger
)

func InitLogger(isProd bool) {
    var cfg zap.Config
    if isProd {
        cfg = zap.NewProductionConfig()
        cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // Уровень логирования в проде
    } else {
        cfg = zap.NewDevelopmentConfig()
        cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel) // Более детализированные логи в dev
    }

    var err error
    Logger, err = cfg.Build()
    if err != nil {
        panic("failed to initialize logger: " + err.Error())
    }
}