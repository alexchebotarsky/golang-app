package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/pure-go-server/client/database"
)

func UpdateArticle(db *database.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var article database.Article
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			handleError(w, fmt.Errorf("error decoding request body: %v", err), http.StatusBadRequest, true)
			return
		}

		upd, err := db.UpdateArticle(id, article)
		if err != nil {
			handleError(w, fmt.Errorf("error updating article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(upd); err != nil {
			log.Printf("%s: %v", logMsgWriteResponse, err)
		}
	}
}
