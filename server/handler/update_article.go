package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-server/client/database"
	"github.com/goodleby/golang-server/tracing"
)

type ArticleUpdater interface {
	UpdateArticle(ctx context.Context, id string, article database.Article) error
}

func UpdateArticle(articleUpdater ArticleUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := tracing.SpanFromContext(ctx)

		id := chi.URLParam(r, "id")

		span.SetTag("id", id)

		article := &database.Article{}
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			HandleError(ctx, w, fmt.Errorf("error decoding request body: %v", err), http.StatusBadRequest, true)
			return
		}

		if err := articleUpdater.UpdateArticle(ctx, id, *article); err != nil {
			HandleError(ctx, w, fmt.Errorf("error updating article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
