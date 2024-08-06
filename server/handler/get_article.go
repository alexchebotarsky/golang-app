package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-app/client"
	"github.com/goodleby/golang-app/model/article"
	"github.com/goodleby/golang-app/tracing"
)

type ArticleSelector interface {
	SelectArticle(ctx context.Context, id int) (*article.Article, error)
}

func GetArticle(articleSelector ArticleSelector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := tracing.SpanFromContext(ctx)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error converting id to int: %v", err), http.StatusBadRequest, false)
			return
		}

		span.SetTag("id", chi.URLParam(r, "id"))

		article, err := articleSelector.SelectArticle(ctx, id)
		if err != nil {
			switch err.(type) {
			case *client.ErrNotFound:
				HandleError(ctx, w, fmt.Errorf("error selecting article: not found: %v", err), http.StatusNotFound, false)
			default:
				HandleError(ctx, w, fmt.Errorf("error selecting article: %v", err), http.StatusInternalServerError, true)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(article)
		handleWritingErr(err)
	}
}
