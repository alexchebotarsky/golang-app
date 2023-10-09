package middleware

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/goodleby/golang-server/metrics"
	"github.com/goodleby/golang-server/processor/event"
)

func Metrics(id string, next event.Handler) event.Handler {
	return func(ctx context.Context, msg *pubsub.Message) {
		start := time.Now()
		next(ctx, msg)
		duration := time.Since(start)

		metrics.RecordEventProcessed(id)
		metrics.ObserveEventDuration(duration.Seconds())
	}
}
