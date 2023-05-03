package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ArticleRemover is an interface that removes an article.
type ArticleRemover interface {
	RemoveArticle(id string) error
}

// RemoveArticle is a handler that removes an article.
func RemoveArticle(articleRemover ArticleRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if err := articleRemover.RemoveArticle(id); err != nil {
			handleError(w, fmt.Errorf("error deleting article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
