package logger

import (
	"log"

	"go.uber.org/zap"
)

var zapLogger *zap.Logger

func InitLogger() {
	var err error
	zapLogger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer zapLogger.Sync()
}

func Info(message string, fields ...zap.Field) {
	zapLogger.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	zapLogger.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	zapLogger.Error(message, fields...)
}
