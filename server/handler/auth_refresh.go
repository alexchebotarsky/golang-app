package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type TokenRefresher interface {
	RefreshToken(ctx context.Context, tokenString string) (token string, expires time.Time, err error)
}

func AuthRefresh(tokenRefresher TokenRefresher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error reading auth token cookie: %v", err), http.StatusUnauthorized, true)
			return
		}

		if time.Now().After(tokenCookie.Expires) {
			HandleError(ctx, w, errors.New("cookie has expired"), http.StatusUnauthorized, true)
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
