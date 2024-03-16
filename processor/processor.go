package processor

import (
	"context"
	"log/slog"

	"cloud.google.com/go/pubsub"
	"github.com/goodleby/golang-app/processor/event"
	"github.com/goodleby/golang-app/processor/middleware"
)

type PubSubClient interface {
	Subscription(id string) *pubsub.Subscription
}

type DBClient interface {
	event.ArticleInserter
}

type Processor struct {
	Events      []event.Event
	Middlewares []middleware.Middleware
	PubSub      PubSubClient
	DB          DBClient
}

func New(ctx context.Context, projectID string, pubsub PubSubClient, db DBClient) (*Processor, error) {
	p := &Processor{}

	p.PubSub = pubsub
	p.DB = db

	p.setupEvents()
	p.setupMiddlewares()

	return p, nil
}

func (p *Processor) Start(ctx context.Context) {
	slog.Info("Processor has started listening to events")
	for _, event := range p.Events {
		go event.Listen(ctx)
	}
}

func (p *Processor) Stop(ctx context.Context) error {
	slog.Debug("Processor has stopped gracefully")

	return nil
}

func (p *Processor) handle(event event.Event) {
	p.Events = append(p.Events, event)
}

func (p *Processor) use(middlewares ...middleware.Middleware) {
	p.Middlewares = append(middlewares, middlewares...)
}

func (p *Processor) setupMiddlewares() {
	for _, event := range p.Events {
		for _, middleware := range p.Middlewares {
			event.Handler = middleware(event.Name, event.Handler)
		}
	}
}

func (p *Processor) setupEvents() {
	p.use(middleware.Trace, middleware.Metrics)

	p.handle(event.Event{
		Name:         "AddArticle",
		Subscription: p.PubSub.Subscription("golang-app-add-article-sub"),
		Handler:      event.AddArticle(p.DB),
	})
}
