package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-app/article"
)

type ArticleInserter interface {
	InsertArticle(ctx context.Context, payload article.Payload) error
}

func AddArticle(articleInserter ArticleInserter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var payload article.Payload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			HandleError(ctx, w, fmt.Errorf("error decoding article payload: %v", err), http.StatusBadRequest, true)
			return
		}

		if err := payload.Validate(); err != nil {
			HandleError(ctx, w, fmt.Errorf("error invalid article payload: %v", err), http.StatusBadRequest, true)
			return
		}

		if err := articleInserter.InsertArticle(ctx, payload); err != nil {
			HandleError(ctx, w, fmt.Errorf("error inserting article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
