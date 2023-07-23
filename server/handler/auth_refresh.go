package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// TokenRefresher is an interface for refreshing JWT.
type TokenRefresher interface {
	RefreshToken(ctx context.Context, tokenString string) (token string, expires time.Time, err error)
}

// AuthRefresh is a handler that creates jwt auth token and stores it in cookie for
// future authorized requests.
func AuthRefresh(tokenRefresher TokenRefresher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error reading auth token cookie: %v", err), http.StatusUnauthorized, true)
			return
		}

		token, expires, err := tokenRefresher.RefreshToken(ctx, tokenCookie.Value)
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error refreshing token: %v", err), http.StatusUnauthorized, true)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  expires,
			HttpOnly: true,
			Path:     "/",
		})

		w.WriteHeader(http.StatusNoContent)
	}
}
