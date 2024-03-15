package event

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

// type ArticleInserter interface {
// 	InsertArticle(ctx context.Context, payload article.Payload) (*article.Article, error)
// }

func AddArticle() Handler {
	return func(ctx context.Context, msg *pubsub.Message) {
		// var payload article.Payload
		// if err := json.Unmarshal(msg.Data, &payload); err != nil {
		// 	msg.Ack()
		// 	HandleError(ctx, fmt.Errorf("error decoding message data: %v", err), true)
		// 	return
		// }

		// _, err := articleInserter.InsertArticle(ctx, payload)
		// if err != nil {
		// 	msg.Nack()
		// 	HandleError(ctx, fmt.Errorf("error adding an article: %v", err), true)
		// 	return
		// }

		log.Printf("Received a pubsub message: %v", msg.Data)

		msg.Ack()
	}
}
