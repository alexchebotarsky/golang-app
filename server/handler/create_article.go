package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/goodleby/pure-go-server/client/database"
)

func CreateArticle(db *database.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var article database.Article
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			handleError(w, fmt.Errorf("error decoding request body: %v", err), http.StatusBadRequest, true)
			return
		}

		article, err := db.CreateArticle(article)
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
