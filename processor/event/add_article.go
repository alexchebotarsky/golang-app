package event

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/goodleby/golang-app/article"
)

type ArticleInserter interface {
	InsertArticle(ctx context.Context, payload article.Payload) error
}

func AddArticle(articleInserter ArticleInserter) Handler {
	return func(ctx context.Context, msg *pubsub.Message) {
		var payload article.Payload
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			msg.Ack()
			HandleError(ctx, fmt.Errorf("error decoding message data: %v", err), true)
			return
		}

		if err := articleInserter.InsertArticle(ctx, payload); err != nil {
			msg.Nack()
			HandleError(ctx, fmt.Errorf("error inserting article: %v", err), true)
			return
		}

		msg.Ack()
	}
}
