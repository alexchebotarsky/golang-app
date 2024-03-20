package tracing

import (
	"fmt"
	"net/http"

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

	span.SetTag("http.method", req.Method)
	span.SetTag("http.url", req.URL.String())

	res, err := tt.RoundTripper.RoundTrip(req)

	span.SetTag("http.status_code", fmt.Sprint(res.StatusCode))

	if err != nil || res.StatusCode >= 400 {
		err = fmt.Errorf("error making request to %s, status code %d: %v", req.URL, res.StatusCode, err)
		span.RecordError(err)
		return nil, err
	}

	return res, err
}
