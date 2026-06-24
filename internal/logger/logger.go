package logger

import (
	"log/slog"
	"os"
)

func New(env string) *slog.Logger {
	options := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if env == "local" {
		return slog.New(slog.NewTextHandler(os.Stdout, options))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, options))
}
