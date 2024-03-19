package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-app/client"
	"github.com/goodleby/golang-app/tracing"
)

type ArticleDeleter interface {
	DeleteArticle(ctx context.Context, id int) error
}

func DeleteArticle(articleDeleter ArticleDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := tracing.SpanFromContext(ctx)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error converting id to int: %v", err), http.StatusBadRequest, false)
			return
		}

		span.SetTag("id", chi.URLParam(r, "id"))

		err = articleDeleter.DeleteArticle(ctx, id)
		if err != nil {
			switch err.(type) {
			case client.ErrNotFound:
				HandleError(ctx, w, fmt.Errorf("error deleting article with id %d: %v", id, err), http.StatusNotFound, false)
			default:
				HandleError(ctx, w, fmt.Errorf("error deleting article: %v", err), http.StatusInternalServerError, true)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
