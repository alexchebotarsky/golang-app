package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-server/client/database"
)

// ArticleUpdater is an interface that updates an article.
type ArticleUpdater interface {
	UpdateArticle(ctx context.Context, id string, article database.Article) error
}

// UpdateArticle is a handler that updates an article.
func UpdateArticle(articleUpdater ArticleUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var article database.Article
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			HandleError(w, fmt.Errorf("error decoding request body: %v", err), http.StatusBadRequest, true)
			return
		}

		if err := articleUpdater.UpdateArticle(r.Context(), id, article); err != nil {
			HandleError(w, fmt.Errorf("error updating article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
