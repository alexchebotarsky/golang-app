package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/goodleby/golang-server/tracing"
)

type Client struct {
	*pubsub.Client
}

func New(ctx context.Context, projectID string) (*Client, error) {
	c := &Client{}
	var err error

	c.Client, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("error creating new pubsub client: %v", err)
	}

	return c, nil
}

func (ps *Client) send(ctx context.Context, topicID string, data []byte) error {
	ctx, span := tracing.StartSpan(ctx, topicID)
	defer span.End()

	topic := ps.Topic(topicID)
	defer topic.Stop()

	var results []*pubsub.PublishResult
	res := topic.Publish(ctx, &pubsub.Message{Data: data, Attributes: tracing.NewCarrier(ctx)})
	results = append(results, res)

	for _, r := range results {
		if _, err := r.Get(ctx); err != nil {
			return fmt.Errorf("error publishing to pubsub: %v", err)
		}
	}
	return nil
}
