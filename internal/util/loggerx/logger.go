package loggerx

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/akasrt/filensy/internal/config/env"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *slog.Logger
	once   sync.Once
)

func Init() *slog.Logger {
	once.Do(func() {
		_ = os.MkdirAll("logs", 0755)

		lj := &lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     90,
			Compress:   true,
		}

		handler := slog.NewJSONHandler(lj, &slog.HandlerOptions{
			Level: parseLevel(),
		})

		logger = slog.New(handler)
	})

	return logger
}

func Get() *slog.Logger {
	if logger == nil {
		panic("logger not initialized")
	}
	return logger
}

func parseLevel() slog.Level {
	level := env.GetEnv(env.LogLevel)
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
