package middleware

import (
	"context"

	"github.com/goodleby/golang-app/processor/event"
	"github.com/goodleby/golang-app/tracing"
)

func Trace(name string, next event.Handler) event.Handler {
	return func(ctx context.Context, msg *event.Message) {
		ctx, span := tracing.StartSpanFromCarrier(ctx, msg.Attributes, name)
		defer span.End()

		span.SetTag("event.name", name)

		next(ctx, msg)

		span.SetTag("event.status", msg.Status)
	}
}
