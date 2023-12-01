package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-app/article"
)

type AddArticlePublisher interface {
	PublishAddArticle(ctx context.Context, payload article.Payload) error
}

func AddArticlePubSub(publisher AddArticlePublisher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var payload article.Payload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			HandleError(ctx, w, fmt.Errorf("error decoding article payload: %v", err), http.StatusBadRequest, false)
			return
		}

		if err := payload.Validate(); err != nil {
			HandleError(ctx, w, fmt.Errorf("error invalid article payload: %v", err), http.StatusBadRequest, false)
			return
		}

		err := publisher.PublishAddArticle(ctx, payload)
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error publishing add article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
