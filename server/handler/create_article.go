package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goodleby/pure-go-server/client/database"
)

// ArticleCreator is an interface that creates an article.
type ArticleCreator interface {
	CreateArticle(ctx context.Context, article database.Article) error
}

// CreateArticle is a handler that creates an article.
func CreateArticle(articleCreator ArticleCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var article database.Article
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			handleError(w, fmt.Errorf("error decoding request body: %v", err), http.StatusBadRequest, true)
			return
		}

		if err := articleCreator.CreateArticle(r.Context(), article); err != nil {
			handleError(w, fmt.Errorf("error creating article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
