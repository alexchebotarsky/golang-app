package event

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

type Handler = func(ctx context.Context, msg *pubsub.Message)

type Event struct {
	ID           string
	Subscription *pubsub.Subscription
	Handler      Handler
}

func (e *Event) Listen(ctx context.Context) {
	if err := e.Subscription.Receive(ctx, e.Handler); err != nil {
		log.Printf("Error listening to %q subscription: %v", e.ID, err)
	}
}
