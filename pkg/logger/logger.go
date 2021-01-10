package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger   *zap.Logger
	sugar    *zap.SugaredLogger
	cfg      zap.Config
	logLevel zapcore.Level
)

func GetInstance() *zap.SugaredLogger {
	cfg = zap.NewProductionConfig()

	switch level := viper.GetString("log-level"); level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	cfg.Level = zap.NewAtomicLevelAt(logLevel)

	logger, _ := cfg.Build()
	defer logger.Sync() // flushes buffer, if any

	sugar = logger.Sugar()

	return sugar
}
