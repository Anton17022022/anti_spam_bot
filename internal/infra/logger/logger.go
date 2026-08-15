package logger

import (
	"log/slog"
	"os"
)

// New возвращает настроенный структурированный логгер, пишущий JSON в stdout.
func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
