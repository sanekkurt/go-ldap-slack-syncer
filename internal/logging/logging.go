package logging

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggingCtxKey struct{}

var (
	log *zap.SugaredLogger
)

func GetLogger() *zap.SugaredLogger {
	return log
}

func Configure(debugMode bool) (*zap.SugaredLogger, error) { //nolint
	var (
		err    error
		logger *zap.Logger

		logLevel = zap.DebugLevel
	)

	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey: "ts",
			//LevelKey:      "L",
			LevelKey: "level",
			//NameKey:       "N",
			CallerKey:   "",
			FunctionKey: zapcore.OmitKey,
			MessageKey:  "message",
			//StacktraceKey: "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if debugMode {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zapConfig.Development = true
		zapConfig.EncoderConfig.CallerKey = "caller"
		zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	if logger, err = zapConfig.Build(); err != nil {
		return nil, fmt.Errorf("cannot create logger: %w", err)
	}

	log = logger.Sugar()

	return log, nil
}

func ConfigureForTests() error {
	var (
		err      error
		logger   *zap.Logger
		logLevel = zap.InfoLevel
	)

	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey: "",
			//LevelKey:      "L",
			LevelKey: "level",
			//NameKey:       "N",
			CallerKey:   "",
			FunctionKey: zapcore.OmitKey,
			MessageKey:  "message",
			//StacktraceKey: "S",
			LineEnding:  zapcore.DefaultLineEnding,
			EncodeLevel: zapcore.CapitalLevelEncoder,
			//EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
	}

	if logger, err = zapConfig.Build(); err != nil {
		return fmt.Errorf("cannot create logger: %w", err)
	}

	log = logger.Sugar()

	return nil
}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggingCtxKey{}, logger)
}

func GetLoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggingCtxKey{}).(*zap.SugaredLogger); ok {
		return logger
	}

	return log
}
