package event

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/goodleby/golang-app/article"
)

type ArticleInserter interface {
	InsertArticle(ctx context.Context, article article.Article) error
}

func AddArticle(articleInserter ArticleInserter) Handler {
	return func(ctx context.Context, msg *pubsub.Message) {
		article := &article.Article{}
		if err := json.Unmarshal(msg.Data, article); err != nil {
			msg.Ack()
			HandleError(ctx, fmt.Errorf("error decoding message data: %v", err), true)
			return
		}

		if err := articleInserter.InsertArticle(ctx, *article); err != nil {
			msg.Nack()
			HandleError(ctx, fmt.Errorf("error inserting article: %v", err), true)
			return
		}

		msg.Ack()
	}
}
