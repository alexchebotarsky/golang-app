package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/goodleby/pure-go-server/client/database"
)

// ArticleCreator is an interface that creates an article.
type ArticleCreator interface {
	CreateArticle(article database.Article) (database.Article, error)
}

// CreateArticle is a handler that creates an article.
func CreateArticle(articleCreator ArticleCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var article database.Article
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			handleError(w, fmt.Errorf("error decoding request body: %v", err), http.StatusBadRequest, true)
			return
		}

		article, err := articleCreator.CreateArticle(article)
		if err != nil {
			handleError(w, fmt.Errorf("error creating article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(article); err != nil {
			log.Printf("%s: %v", logMsgWriteResponse, err)
		}
	}
}
