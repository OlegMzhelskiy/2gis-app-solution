package log

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func InitializeLogger() {
	var err error

	logger, err = zap.NewDevelopment()
	// logger, err = zap.NewProduction()

	if err != nil {
		panic("Failed to initialize zap logger: " + err.Error())
	}
}

func GetLogger() *zap.Logger {
	return logger
}

func Error(msg string, err error) {
	logger.Error(msg, zap.Error(err))
}

func Fatal(msg string, err error) {
	logger.Fatal(msg, zap.Error(err))
}

func Warning(msg string) {
	logger.Warn(msg)
}

func Info(msg string) {
	logger.Info(msg)
}

func WithField(key string, value any) *zap.Logger {
	return logger.With(zap.Any(key, value))
}

func WithFields(f map[string]any) *zap.Logger {
	fields := make([]zap.Field, 0, len(f))

	for k, v := range f {
		fields = append(fields, zap.Any(k, v))
	}

	return logger.With(fields...)
}
