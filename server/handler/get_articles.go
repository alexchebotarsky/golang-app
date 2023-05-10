package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/goodleby/pure-go-server/client/database"
)

// ArticlesFetcher is an interface that fetches articles.
type AllArticlesFetcher interface {
	FetchAllArticles(ctx context.Context) ([]database.Article, error)
}

// GetArticles is a handler that fetches articles.
func GetAllArticles(articleFetcher AllArticlesFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := articleFetcher.FetchAllArticles(r.Context())
		if err != nil {
			handleError(w, fmt.Errorf("error fetching articles: %v", err), http.StatusInternalServerError, true)
			return
		}

		if len(articles) == 0 {
			handleError(w, fmt.Errorf("articles not found"), http.StatusNotFound, false)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(articles); err != nil {
			log.Printf("%s: %v", logMsgWriteResponse, err)
		}
	}
}
