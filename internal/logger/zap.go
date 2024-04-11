package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(defaultLogLevel zapcore.Level) (logger *Logger, err error) {
	logger = &Logger{}
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := logger.getLoggerTee(config, defaultLogLevel)
	logger.Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.FatalLevel))

	return
}

func (logger Logger) getLoggerTee(config zapcore.EncoderConfig, defaultLogLevel zapcore.Level) (core zapcore.Core) {
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	core = zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	return
}
