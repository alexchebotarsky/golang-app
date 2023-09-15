package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-server/client/database"
)

type ArticleInserter interface {
	InsertArticle(ctx context.Context, article database.Article) error
}

func AddArticle(articleInserter ArticleInserter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var article database.Article
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			HandleError(ctx, w, fmt.Errorf("error decoding request body: %v", err), http.StatusBadRequest, true)
			return
		}

		if err := articleInserter.InsertArticle(ctx, article); err != nil {
			HandleError(ctx, w, fmt.Errorf("error inserting article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}