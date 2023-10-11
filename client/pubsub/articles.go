package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goodleby/golang-app/article"
)

func (c *Client) PublishAddArticle(ctx context.Context, article *article.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("error marshalling article: %v", err)
	}

	c.send(ctx, "golang-app-add-article", data)

	return nil
}
