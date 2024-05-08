package logger

import (
	"os"

	"log/slog"
)

func Init(level slog.Level, format string) {
	opts := &slog.HandlerOptions{
		ReplaceAttr: renameMessageKey,
		Level:       level,
	}

	var handler slog.Handler
	switch format {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	case "json":
		fallthrough
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)
}

func renameMessageKey(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.MessageKey {
		return slog.Attr{
			Key:   MessageKey,
			Value: a.Value,
		}
	}

	return a
}

const MessageKey = "message"
