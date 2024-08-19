package log

import (
	"log/slog"
	"os"
)

func Configure() {
	level := getLogLevel()
	logConfig := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	slog.SetDefault(logConfig)

	slog.Debug("Log configured")
}

func getLogLevel() slog.Level {
	level := os.Getenv("LOG_LEVEL")

	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
