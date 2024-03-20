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

func (tt *TracedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	_, span := StartSpan(r.Context(), "RoundTrip")
	defer span.End()

	span.SetTag("http.method", r.Method)
	span.SetTag("http.url", r.URL.String())

	res, err := tt.RoundTripper.RoundTrip(r)

	span.SetTag("http.status_code", fmt.Sprint(res.StatusCode))

	if err != nil || res.StatusCode >= 400 {
		err = fmt.Errorf("error making request to %s, status code %d: %v", r.URL, res.StatusCode, err)
		span.RecordError(err)
		return nil, err
	}

	return res, err
}
