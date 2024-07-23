package middleware

import (
	"context"
	"time"

	"github.com/goodleby/golang-app/metrics"
	"github.com/goodleby/golang-app/processor/event"
)

func Metrics(eventName string, next event.Handler) event.Handler {
	return func(ctx context.Context, msg *event.Message) {
		start := time.Now()
		next(ctx, msg)
		duration := time.Since(start)

		metrics.RecordEventProcessed(eventName, msg.Status)
		metrics.ObserveEventDuration(eventName, duration)
	}
}
