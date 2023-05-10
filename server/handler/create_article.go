package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goodleby/pure-go-server/client/database"
)

// ArticleAdder is an interface that adds an article.
type ArticleAdder interface {
	AddArticle(ctx context.Context, article database.Article) error
}

// AddArticle is a handler that adds an article.
func AddArticle(articleAdder ArticleAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var article database.Article
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			handleError(w, fmt.Errorf("error decoding request body: %v", err), http.StatusBadRequest, true)
			return
		}

		if err := articleAdder.AddArticle(r.Context(), article); err != nil {
			handleError(w, fmt.Errorf("error adding article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
