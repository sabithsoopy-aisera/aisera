package aisera

import (
	"os"

	"golang.org/x/exp/slog"
)

var logger *slog.Logger

func init() {
	logHandler := slog.HandlerOptions{Level: slog.LevelInfo}
	if os.Getenv("DEBUG") != "" {
		logHandler.Level = slog.LevelDebug
	}
	logger = slog.New(logHandler.NewJSONHandler(os.Stderr)).WithGroup("aisera")
}

func SetLogger(l *slog.Logger) {
	logger = l
}
