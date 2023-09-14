package applog

import (
	"log/slog"
	"os"
)

func NewDefaultLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	return logger
}
