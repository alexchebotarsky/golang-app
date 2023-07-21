package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ArticleRemover is an interface that removes an article.
type ArticleRemover interface {
	RemoveArticle(ctx context.Context, id string) error
}

// RemoveArticle is a handler that removes an article.
func RemoveArticle(articleRemover ArticleRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		id := chi.URLParam(r, "id")

		span.SetAttributes(attribute.String("id", id))

		if err := articleRemover.RemoveArticle(ctx, id); err != nil {
			HandleError(ctx, w, fmt.Errorf("error removing article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
