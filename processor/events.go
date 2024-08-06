package processor

import (
	"github.com/goodleby/golang-app/processor/event"
	"github.com/goodleby/golang-app/processor/handler"
	"github.com/goodleby/golang-app/processor/middleware"
)

func (p *Processor) setupEvents() {
	p.use(middleware.Trace, middleware.Metrics)

	p.handle(event.Event{
		Name:           "AddArticle",
		SubscriptionID: "golang-app-add-article-sub",
		Handler:        handler.AddArticle(p.Clients.DB),
		Throttle:       1,
	})
}
