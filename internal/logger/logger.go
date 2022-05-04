package logger

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func New(level string, file string) (*Logger, error) {
	var zapLevel zap.AtomicLevel

	switch level {
	case "error":
		zapLevel = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "warn":
		zapLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "info":
		zapLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "debug":
		zapLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	default:
		return nil, errors.New("unsupported level")
	}

	config := zap.NewProductionConfig()

	config.Level = zapLevel
	config.Encoding = "json"

	if file != "stdout" && file != "stderr" {
		config.OutputPaths = append(config.OutputPaths, file)
	}

	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{logger: logger}, nil
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}
