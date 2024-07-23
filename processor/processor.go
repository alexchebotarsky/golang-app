package processor

import (
	"context"
	"fmt"
	"log/slog"

	"cloud.google.com/go/pubsub"
	"github.com/goodleby/golang-app/processor/event"
	"github.com/goodleby/golang-app/processor/handler"
	"github.com/goodleby/golang-app/processor/middleware"
)

type Processor struct {
	Events      []event.Event
	Middlewares []event.Middleware
	Clients     Clients
}

type Clients struct {
	PubSub PubSubClient
	DB     DBClient
}

type PubSubClient interface {
	Subscription(id string) *pubsub.Subscription
}

type DBClient interface {
	handler.ArticleInserter
}

func New(ctx context.Context, clients Clients) (*Processor, error) {
	var p Processor

	p.Clients = clients

	p.setupEvents()

	return &p, nil
}

func (p *Processor) Start(ctx context.Context, errc chan<- error) {
	for _, e := range p.Events {
		// Gather global processor middlewares and event specific middlewares.
		middlewares := make([]event.Middleware, 0, len(p.Middlewares)+len(e.Middlewares))
		middlewares = append(middlewares, p.Middlewares...)
		middlewares = append(middlewares, e.Middlewares...)

		// Apply relevant middlewares before listening to the event.
		for _, middleware := range middlewares {
			e.Handler = middleware(e.Name, e.Handler)
		}

		go e.Listen(ctx, errc)
	}

	slog.Info(fmt.Sprintf("PubSub event processor listening to %d events", len(p.Events)))
}

func (p *Processor) Stop(ctx context.Context) error {
	return nil
}

func (p *Processor) handle(e event.Event) {
	e.Subscription = p.Clients.PubSub.Subscription(e.SubscriptionID)
	e.Subscription.ReceiveSettings.MaxOutstandingMessages = e.Throttle

	p.Events = append(p.Events, e)
}

func (p *Processor) use(middlewares ...event.Middleware) {
	p.Middlewares = append(p.Middlewares, middlewares...)
}

func (p *Processor) setupEvents() {
	p.use(middleware.Trace, middleware.Metrics)

	p.handle(event.Event{
		Name:           "AddArticle",
		SubscriptionID: "golang-app-add-article-sub",
		Handler:        handler.AddArticle(p.Clients.DB),
		Throttle:       1,
	})
}
