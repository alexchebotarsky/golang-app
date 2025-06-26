package event

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type Event struct {
	Name           string
	SubscriptionID string
	Handler        Handler
	Throttle       int

	Subscription *pubsub.Subscription
	Middlewares  []Middleware
}

type Handler = func(ctx context.Context, msg *Message)

type Middleware func(eventName string, next Handler) Handler

func (e *Event) Listen(ctx context.Context, errc chan<- error) {
	if e.Subscription == nil {
		errc <- fmt.Errorf("Subscription %q is nil", e.Name)
		return
	}

	err := e.Subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		e.Handler(ctx, &Message{Message: msg})
	})
	if err != nil {
		errc <- fmt.Errorf("Error listening to %q subscription: %v", e.Name, err)
		return
	}
}

type Message struct {
	*pubsub.Message
	Status string
}

func (m *Message) Ack() {
	if m.Status == "" {
		m.Status = StatusOK
	}
	m.Message.Ack()
}

func (m *Message) AckWithResult() *pubsub.AckResult {
	if m.Status == "" {
		m.Status = StatusOK
	}
	return m.Message.AckWithResult()
}

func (m *Message) Nack() {
	if m.Status == "" {
		m.Status = StatusRetry
	}
	m.Message.Nack()
}

func (m *Message) NackWithResult() *pubsub.AckResult {
	if m.Status == "" {
		m.Status = StatusRetry
	}
	return m.Message.NackWithResult()
}

func (m *Message) SetStatus(status string) {
	m.Status = status
}

const (
	StatusOK     = "OK"
	StatusFailed = "Failed"
	StatusRetry  = "Retry"
)
