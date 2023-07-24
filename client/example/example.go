package example

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/goodleby/golang-server/client"
	"github.com/goodleby/golang-server/config"
	"github.com/goodleby/golang-server/tracing"
)

// Client is an example client.
type Client struct {
	ExampleEndpoint string
	HTTPClient      *http.Client
}

// New creates a new example client.
func New(config *config.Config) (*Client, error) {
	var c Client

	c.ExampleEndpoint = config.ExampleEndpoint

	c.HTTPClient = client.NewHTTPClient(client.Parameters{
		Timeout: 3 * time.Second,
	})

	return &c, nil
}

// ExampleData is the data we expect to receive from example endpoint.
type ExampleData struct {
	Name      string `json:"name"`
	Height    string `json:"height"`
	Mass      string `json:"mass"`
	BirthYear string `json:"birth_year"`
	Gender    string `json:"gender"`
}

// FetchExampleData sends get request to the example endpoint.
func (c *Client) FetchExampleData(ctx context.Context) (*ExampleData, error) {
	ctx, span := tracing.Span(ctx, "FetchExampleData")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.ExampleEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %v", err)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing http request: %v", err)
	}

	var exampleData ExampleData
	if err := json.NewDecoder(res.Body).Decode(&exampleData); err != nil {
		return nil, fmt.Errorf("error decoding request body as json: %v", err)
	}

	return &exampleData, nil
}