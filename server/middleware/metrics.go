package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-server/metrics"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		crw := customResponseWriter{ResponseWriter: w}
		start := time.Now()

		next.ServeHTTP(&crw, r)

		duration := time.Since(start)

		routeID := fmt.Sprintf("%s %s", r.Method, chi.RouteContext(r.Context()).RoutePattern())

		metrics.ObserveTimeToProcess(routeID, duration.Seconds())
		metrics.RecordRequestStatusCode(crw.status, routeID)
	})
}

type customResponseWriter struct {
	http.ResponseWriter
	status int
}

func (crw *customResponseWriter) WriteHeader(status int) {
	crw.status = status
	crw.ResponseWriter.WriteHeader(status)
}
