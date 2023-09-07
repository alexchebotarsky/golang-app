package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthLoginPayload struct {
	Role string `json:"role"`
	Key  string `json:"key"`
}

type TokenCreator interface {
	NewToken(ctx context.Context, role, key string) (token string, expires time.Time, err error)
}

func AuthLogin(tokenCreator TokenCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var payload AuthLoginPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			HandleError(ctx, w, fmt.Errorf("error decoding auth payload: %v", err), http.StatusBadRequest, true)
			return
		}

		token, expires, err := tokenCreator.NewToken(ctx, payload.Role, payload.Key)
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error creating new auth token: %v", err), http.StatusUnauthorized, true)
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
