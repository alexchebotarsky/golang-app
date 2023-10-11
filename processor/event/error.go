package event

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/goodleby/golang-app/tracing"
)

func HandleError(ctx context.Context, err error, shouldLog bool) {
	span := tracing.SpanFromContext(ctx)

	span.RecordError(err)

	if shouldLog {
		slog.Error(fmt.Sprintf("Event error: %v", err))
	}
}
