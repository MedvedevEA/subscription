package logger

import (
	"log"
	"subscription/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func MustNew(cfg *config.Log) *zap.SugaredLogger {
	atomicLevel, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		log.Fatalf("failed to parse level: %v\n", err)
	}

	logger, err := zap.Config{
		Level: atomicLevel,

		Encoding:    "json",
		OutputPaths: cfg.OutPutPaths,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "timestamp",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
		DisableStacktrace: true,
	}.Build()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v\n", err)
	}

	return logger.Sugar()
}
