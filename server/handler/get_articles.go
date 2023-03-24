package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/goodleby/pure-go-server/client/database"
)

func GetArticles(db *database.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := db.FetchArticles()
		if err != nil {
			handleError(w, fmt.Errorf("error fetching articles: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(articles); err != nil {
			log.Printf("%s: %v", logMsgWriteResponse, err)
		}
	}
}
