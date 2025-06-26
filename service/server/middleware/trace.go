package middleware

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-app/tracing"
)

func Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracing.StartSpan(r.Context(), "UnknownServerHandler")
		defer span.End()

		span.SetTag("http.method", r.Method)
		span.SetTag("http.url", r.URL.String())

		next.ServeHTTP(w, r.WithContext(ctx))

		routeID := fmt.Sprintf("%s %s", r.Method, chi.RouteContext(ctx).RoutePattern())
		span.SetName(routeID)
	})
}
