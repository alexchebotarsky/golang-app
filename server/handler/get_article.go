package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-server/client/database"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ArticleFetcher is an interface that fetches an article.
type ArticleFetcher interface {
	FetchArticle(ctx context.Context, id string) (*database.Article, error)
}

// GetArticle is a handler that fetches an article.
func GetArticle(articleFetcher ArticleFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		id := chi.URLParam(r, "id")

		span.SetAttributes(attribute.String("id", id))

		article, err := articleFetcher.FetchArticle(ctx, id)
		if err != nil {
			switch err.(type) {
			case *database.ErrNotFound:
				HandleError(ctx, w, fmt.Errorf("article with id %q not found: %v", id, err), http.StatusNotFound, false)
			default:
				HandleError(ctx, w, fmt.Errorf("error fetching article: %v", err), http.StatusInternalServerError, true)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(article)
		handleWritingErr(err)
	}
}
