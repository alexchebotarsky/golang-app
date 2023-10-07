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

type ArticleSelector interface {
	SelectArticle(ctx context.Context, id string) (*database.Article, error)
}

func GetArticle(articleSelector ArticleSelector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := tracing.SpanFromContext(ctx)

		id := chi.URLParam(r, "id")

		span.SetTag("id", id)

		article, err := articleSelector.SelectArticle(ctx, id)
		if err != nil {
			switch err.(type) {
			case *database.ErrNotFound:
				HandleError(ctx, w, fmt.Errorf("article with id %q not found: %v", id, err), http.StatusNotFound, false)
			default:
				HandleError(ctx, w, fmt.Errorf("error selecting article: %v", err), http.StatusInternalServerError, true)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(article)
		handleWritingErr(err)
	}
}
