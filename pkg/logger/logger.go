package logger

import (
	"log"

	"go.uber.org/zap"
)

// New inisialisasi logger Zap yang terstruktur.
func New() *zap.Logger {
	// NewProduction memberikan logger yang siap untuk production:
	// - Level: Info
	// - Output: JSON
	// - Performa tinggi
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	return logger
}
