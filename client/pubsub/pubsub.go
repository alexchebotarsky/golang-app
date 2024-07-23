package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/goodleby/golang-app/tracing"
)

type Client struct {
	*pubsub.Client
}

func New(ctx context.Context, projectID string) (*Client, error) {
	var c Client
	var err error

	c.Client, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("error creating new pubsub client: %v", err)
	}

	return &c, nil
}

func (c *Client) Subscription(id string) *pubsub.Subscription {
	return c.Client.Subscription(id)
}

func (c *Client) send(ctx context.Context, topicID string, data []byte) error {
	ctx, span := tracing.StartSpan(ctx, topicID)
	defer span.End()

	topic := c.Client.Topic(topicID)
	defer topic.Stop()

	result := topic.Publish(ctx, &pubsub.Message{Data: data, Attributes: tracing.NewCarrier(ctx)})
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish message %s to topic %q: %v", id, topicID, err)
	}

	return nil
}
