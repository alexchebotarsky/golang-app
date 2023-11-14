package middleware

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/goodleby/golang-app/metrics"
	"github.com/goodleby/golang-app/processor/event"
)

func Metrics(name string, next event.Handler) event.Handler {
	return func(ctx context.Context, msg *pubsub.Message) {
		start := time.Now()
		next(ctx, msg)
		duration := time.Since(start)

		metrics.RecordEventProcessed(name)
		metrics.ObserveEventDuration(duration.Seconds())
	}
}
