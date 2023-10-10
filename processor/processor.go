package processor

import (
	"context"

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

func (p *Processor) Run(ctx context.Context) {
	go p.listenToEvents(ctx)
}

func (p *Processor) listenToEvents(ctx context.Context) {
	for _, event := range p.Events {
		go event.Listen(ctx)
	}
}

func (p *Processor) handle(subID string, handler event.Handler) {
	sub := event.Event{
		ID:           subID,
		Subscription: p.PubSub.Subscription(subID),
		Handler:      handler,
	}
	p.Events = append(p.Events, sub)
}

func (p *Processor) use(middlewares ...middleware.Middleware) {
	p.Middlewares = append(middlewares, middlewares...)
}

func (p *Processor) setupMiddlewares() {
	for _, event := range p.Events {
		for _, middleware := range p.Middlewares {
			event.Handler = middleware(event.ID, event.Handler)
		}
	}
}

func (p *Processor) setupEvents() {
	p.use(middleware.Trace, middleware.Metrics)

	p.handle("golang-server-add-article", event.AddArticle(p.DB))
}
