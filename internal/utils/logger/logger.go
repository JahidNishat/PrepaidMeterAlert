package logger

import (
	"log/slog"
	"os"
)

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// InitLogger configures the default slog logger with the given level and output format.
func InitLogger(level slog.Level, format Format) {
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if format == FormatJSON {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	slog.SetDefault(slog.New(handler))
}
