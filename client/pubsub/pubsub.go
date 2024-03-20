package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/goodleby/golang-app/tracing"
)

type Client struct {
	*pubsub.Client
	envTag string
}

func New(ctx context.Context, projectID, envTag string) (*Client, error) {
	var c Client
	var err error

	c.envTag = envTag

	c.Client, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("error creating new pubsub client: %v", err)
	}

	return &c, nil
}

func (c *Client) Subscription(id string) *pubsub.Subscription {
	return c.Client.Subscription(c.withEnvTag(id))
}

func (c *Client) send(ctx context.Context, topicID string, data []byte) error {
	ctx, span := tracing.StartSpan(ctx, topicID)
	defer span.End()

	topic := c.Client.Topic(c.withEnvTag(topicID))
	defer topic.Stop()

	var results []*pubsub.PublishResult
	res := topic.Publish(ctx, &pubsub.Message{Data: data, Attributes: tracing.NewCarrier(ctx)})
	results = append(results, res)

	for _, r := range results {
		_, err := r.Get(ctx)
		if err != nil {
			return fmt.Errorf("error publishing to pubsub: %v", err)
		}
	}
	return nil
}

func (c *Client) withEnvTag(id string) string {
	if c.envTag == "" || id == "" {
		return id
	}

	return fmt.Sprintf("%s-%s", id, c.envTag)
}
