package event

import (
	"context"
	"log"

	"github.com/goodleby/golang-app/tracing"
)

func HandleError(ctx context.Context, err error, shouldLog bool) {
	span := tracing.SpanFromContext(ctx)

	span.RecordError(err)

	if shouldLog {
		log.Printf("Event error: %v", err)
	}
}
