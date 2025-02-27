package logger

import (
	"log/slog"
	"os"
)

func InitLogger(cfg string) {
	var handler slog.Handler
	switch cfg {
	case "production":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	case "development":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger) // Устанавливаем глобальный логгер
}

func GetLogger(module string) *slog.Logger {
	return slog.Default().With("module", module)
}
