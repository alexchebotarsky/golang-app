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

type TracedTransport struct {
	http.RoundTripper
}

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

func NewTracedTransport(base http.RoundTripper) *TracedTransport {
	return &TracedTransport{base}
}

type Parameters struct {
	Timeout time.Duration
}

func NewHTTPClient(params Parameters) *http.Client {
	c := http.Client{
		Timeout:   params.Timeout,
		Transport: NewTracedTransport(otelhttp.NewTransport(http.DefaultTransport)),
	}

	return &c
}
