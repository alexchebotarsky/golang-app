package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/pure-go-server/client/database"
)

func GetArticle(db *database.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		article, err := db.FetchArticle(id)
		if err != nil {
			handleError(w, fmt.Errorf("error fetching article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(article); err != nil {
			log.Printf("%s: %v", logMsgWriteResponse, err)
		}
	}
}
