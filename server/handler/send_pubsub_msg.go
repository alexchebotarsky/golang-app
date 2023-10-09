package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-server/article"
)

type AddArticlePublisher interface {
	PublishAddArticle(ctx context.Context, article *article.Article) error
}

func SendPubSubMsg(publisher AddArticlePublisher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		article := &article.Article{}
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			HandleError(ctx, w, fmt.Errorf("error decoding request body: %v", err), http.StatusBadRequest, true)
			return
		}

		if err := article.Validate(); err != nil {
			HandleError(ctx, w, fmt.Errorf("error invalid article: %v", err), http.StatusBadRequest, true)
			return
		}

		err := publisher.PublishAddArticle(ctx, article)
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error publishing add article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
