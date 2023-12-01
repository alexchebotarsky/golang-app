package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/goodleby/golang-app/client"
)

type TokenRefresher interface {
	RefreshToken(ctx context.Context, tokenString string) (token string, expires time.Time, err error)
}

func AuthRefresh(tokenRefresher TokenRefresher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error reading auth token cookie: %v", err), http.StatusUnauthorized, false)
			return
		}

		token, expires, err := tokenRefresher.RefreshToken(ctx, tokenCookie.Value)
		if err != nil {
			switch err.(type) {
			case client.ErrUnauthorized:
				HandleError(ctx, w, fmt.Errorf("error refreshing token: %v", err), http.StatusUnauthorized, false)
			default:
				HandleError(ctx, w, fmt.Errorf("error refreshing token: %v", err), http.StatusInternalServerError, true)
			}
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
