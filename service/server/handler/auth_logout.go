package handler

import (
	"net/http"
	"time"
)

// AuthLogout immediately expires auth token cookie on the client.
func AuthLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusNoContent)
}
