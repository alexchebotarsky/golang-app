package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-app/article"
)

type ArticleInserter interface {
	InsertArticle(ctx context.Context, payload article.Payload) (*article.Article, error)
}

func AddArticle(articleInserter ArticleInserter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var payload article.Payload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error decoding article payload: %v", err), http.StatusBadRequest, false)
			return
		}

		err = payload.Validate()
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error invalid article payload: %v", err), http.StatusBadRequest, false)
			return
		}

		article, err := articleInserter.InsertArticle(ctx, payload)
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error adding an article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(article)
		handleWritingErr(err)
	}
}
