package handler

import (
	"net/http"
	"time"
)

// AuthLogout immediately expires auth token cookie on the client.
func AuthLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	w.WriteHeader(http.StatusNoContent)
}
