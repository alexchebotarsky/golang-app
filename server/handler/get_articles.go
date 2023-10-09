package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-server/article"
)

type AllArticlesSelector interface {
	SelectAllArticles(ctx context.Context) ([]article.Article, error)
}

func GetAllArticles(articleSelector AllArticlesSelector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		articles, err := articleSelector.SelectAllArticles(ctx)
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error selecting articles: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(articles)
		handleWritingErr(err)
	}
}
