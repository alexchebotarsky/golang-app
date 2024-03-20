package event

import (
	"context"
	"fmt"
	"log/slog"

	"cloud.google.com/go/pubsub"
)

type Event struct {
	Name         string
	Subscription *pubsub.Subscription
	Handler      Handler
}

type Handler = func(ctx context.Context, msg *pubsub.Message)

func (e *Event) Listen(ctx context.Context) {
	err := e.Subscription.Receive(ctx, e.Handler)
	if err != nil {
		slog.Error(fmt.Sprintf("Error listening to %q subscription: %v", e.Name, err), "subscriptionName", e.Name)
	}
}
