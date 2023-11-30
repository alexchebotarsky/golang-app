package handler

import (
	"context"
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-app/client/database"
	"github.com/goodleby/golang-app/tracing"
)

type ArticleDeleter interface {
	DeleteArticle(ctx context.Context, id string) error
}

func DeleteArticle(articleDeleter ArticleDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := tracing.SpanFromContext(ctx)

		id := chi.URLParam(r, "id")

		span.SetTag("id", id)

		if err := articleDeleter.DeleteArticle(ctx, id); err != nil {
			switch err.(type) {
			case database.ErrNotFound:
				HandleError(ctx, w, fmt.Errorf("error deleting article with id %q: %v", id, err), http.StatusNotFound, false)
			default:
				HandleError(ctx, w, fmt.Errorf("error deleting article: %v", err), http.StatusInternalServerError, true)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
