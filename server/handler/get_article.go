package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-server/client/database"
)

// ArticleFetcher is an interface that fetches an article.
type ArticleFetcher interface {
	FetchArticle(ctx context.Context, id string) (*database.Article, error)
}

// GetArticle is a handler that fetches an article.
func GetArticle(articleFetcher ArticleFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		article, err := articleFetcher.FetchArticle(r.Context(), id)
		if err != nil {
			switch err.(type) {
			case *database.ErrNotFound:
				HandleError(w, fmt.Errorf("article with id %q not found: %v", id, err), http.StatusNotFound, false)
			default:
				HandleError(w, fmt.Errorf("error fetching article: %v", err), http.StatusInternalServerError, true)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(article); err != nil {
			log.Printf("%s: %v", logMsgWriteResponse, err)
		}
	}
}
