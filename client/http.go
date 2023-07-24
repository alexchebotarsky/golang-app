package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goodleby/golang-server/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// TracedTransport implement the http.RoundTripper interface and adds tracing.
type TracedTransport struct {
	http.RoundTripper
}

// RoundTrip executes a single HTTP transaction and is part of the
// http.RoundTripper interface.
func (tt TracedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	_, span := tracing.Span(req.Context(), "RoundTrip")
	defer span.End()

	span.SetAttributes(attribute.String("URL", req.URL.String()))

	res, err := tt.RoundTripper.RoundTrip(req)
	if err != nil || res.StatusCode >= 400 {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error making request to %s, status code: %d", req.URL, res.StatusCode))
	}

	span.SetAttributes(attribute.Int("StatusCode", res.StatusCode))

	return res, err
}

// NewTracedTransport wraps provided base transport with traced transport.
func NewTracedTransport(base http.RoundTripper) *TracedTransport {
	return &TracedTransport{base}
}

// Parameters contains all parameters for the new http client.
type Parameters struct {
	Timeout time.Duration
}

// NewHTTPClient creates an http client with provided parameters and custom
// transport that collects traces and metrics.
func NewHTTPClient(params Parameters) *http.Client {
	c := http.Client{
		Timeout:   params.Timeout,
		Transport: NewTracedTransport(otelhttp.NewTransport(http.DefaultTransport)),
	}

	return &c
}
