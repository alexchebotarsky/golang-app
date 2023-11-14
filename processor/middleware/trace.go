package middleware

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/goodleby/golang-app/processor/event"
	"github.com/goodleby/golang-app/tracing"
)

func Trace(name string, next event.Handler) event.Handler {
	return func(ctx context.Context, msg *pubsub.Message) {
		ctx, span := tracing.StartSpanFromCarrier(ctx, msg.Attributes, name)
		defer span.End()

		span.SetTag("EventName", name)

		next(ctx, msg)
	}
}
