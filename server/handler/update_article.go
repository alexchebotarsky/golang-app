package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-app/article"
	"github.com/goodleby/golang-app/client"
	"github.com/goodleby/golang-app/tracing"
)

type ArticleUpdater interface {
	UpdateArticle(ctx context.Context, id int, payload article.Payload) (*article.Article, error)
}

func UpdateArticle(articleUpdater ArticleUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := tracing.SpanFromContext(ctx)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error converting id to int: %v", err), http.StatusBadRequest, false)
			return
		}

		span.SetTag("id", chi.URLParam(r, "id"))

		var payload article.Payload
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error decoding article payload: %v", err), http.StatusBadRequest, false)
			return
		}

		err = payload.Validate()
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error invalid article payload: %v", err), http.StatusBadRequest, false)
			return
		}

		article, err := articleUpdater.UpdateArticle(ctx, id, payload)
		if err != nil {
			switch err.(type) {
			case *client.ErrNotFound:
				HandleError(ctx, w, fmt.Errorf("error updating article: %v", err), http.StatusNotFound, false)
			default:
				HandleError(ctx, w, fmt.Errorf("error updating article: %v", err), http.StatusInternalServerError, true)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(article)
		handleWritingErr(err)
	}
}
