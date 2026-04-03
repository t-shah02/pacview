package utils

import (
	"log/slog"
	"os"
	"sync"
)

var (
	instance *slog.Logger
	once     sync.Once
)

func GetLogger() *slog.Logger {
	once.Do(func() {
		opts := &slog.HandlerOptions{Level: slog.LevelDebug}
		instance = slog.New(slog.NewTextHandler(os.Stdout, opts))
	})
	return instance
}
