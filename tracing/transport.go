package tracing

import (
	"fmt"
	"net/http"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type TracedTransport struct {
	http.RoundTripper
}

func NewTracedTransport(base http.RoundTripper) *TracedTransport {
	return &TracedTransport{otelhttp.NewTransport(base)}
}

func (tt *TracedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	_, span := StartSpan(req.Context(), "RoundTrip")
	defer span.End()

	span.SetTag("URL", req.URL.String())

	res, err := tt.RoundTripper.RoundTrip(req)
	if err != nil || res.StatusCode >= 400 {
		span.RecordError(fmt.Errorf("error making request to %s, status code %d: %v", req.URL, res.StatusCode, err))
	}

	span.SetTag("StatusCode", strconv.Itoa(res.StatusCode))

	return res, err
}
