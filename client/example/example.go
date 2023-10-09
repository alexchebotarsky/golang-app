package example

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/goodleby/golang-server/tracing"
)

type Client struct {
	ExampleEndpoint string
	HTTPClient      *http.Client
}

func New(endpoint string) (*Client, error) {
	c := &Client{}

	c.ExampleEndpoint = endpoint

	c.HTTPClient = &http.Client{
		Timeout:   3 * time.Second,
		Transport: tracing.NewTracedTransport(http.DefaultTransport),
	}

	return c, nil
}

type ExampleData struct {
	Name      string `json:"name"`
	Height    string `json:"height"`
	Mass      string `json:"mass"`
	BirthYear string `json:"birth_year"`
	Gender    string `json:"gender"`
}

func (c *Client) FetchExampleData(ctx context.Context) (*ExampleData, error) {
	ctx, span := tracing.StartSpan(ctx, "FetchExampleData")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.ExampleEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %v", err)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing http request: %v", err)
	}

	exampleData := &ExampleData{}
	if err := json.NewDecoder(res.Body).Decode(exampleData); err != nil {
		return nil, fmt.Errorf("error decoding request body as json: %v", err)
	}

	return exampleData, nil
}
