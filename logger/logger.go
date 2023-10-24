package logger

import (
	"os"

	"log/slog"
)

func Init() {
	opts := &slog.HandlerOptions{
		ReplaceAttr: renameMessageKey,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
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
