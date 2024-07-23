package handler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/goodleby/golang-app/processor/event"
	"github.com/goodleby/golang-app/tracing"
)

func HandleError(ctx context.Context, msg *event.Message, err error, retry bool) {
	span := tracing.SpanFromContext(ctx)

	span.SetTag("event.status", msg.Status)
	span.RecordError(err)

	if retry {
		msg.SetStatus(event.StatusRetry)
		msg.Nack()
	} else {
		msg.SetStatus(event.StatusFailed)
		msg.Ack()
	}

	slog.Error(fmt.Sprintf("Event error: %v", err))
}
