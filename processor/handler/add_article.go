package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goodleby/golang-app/article"
	"github.com/goodleby/golang-app/processor/event"
)

type ArticleInserter interface {
	InsertArticle(ctx context.Context, payload article.Payload) (*article.Article, error)
}

func AddArticle(articleInserter ArticleInserter) event.Handler {
	return func(ctx context.Context, msg *event.Message) {
		var payload article.Payload
		err := json.Unmarshal(msg.Data, &payload)
		if err != nil {
			HandleError(ctx, msg, fmt.Errorf("error decoding message data: %v", err), false)
			return
		}

		_, err = articleInserter.InsertArticle(ctx, payload)
		if err != nil {
			HandleError(ctx, msg, fmt.Errorf("error adding an article: %v", err), true)
			return
		}

		msg.Ack()
	}
}
