package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/pure-go-server/client/database"
)

func DeleteArticle(db *database.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if err := db.DeleteArticle(id); err != nil {
			handleError(w, fmt.Errorf("error deleting article: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
