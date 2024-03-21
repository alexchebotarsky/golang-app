package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-app/client"
	"github.com/goodleby/golang-app/client/auth"
	"github.com/goodleby/golang-app/server/handler"
)

type TokenAccessReader interface {
	ReadTokenAccess(ctx context.Context, token string) (auth.AccessLevel, error)
}

func Auth(tokenReader TokenAccessReader, expectedAccess auth.AccessLevel) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tokenCookie, err := r.Cookie("token")
			if err != nil {
				handler.HandleError(ctx, w, fmt.Errorf("error reading auth token cookie: %v", err), http.StatusUnauthorized, false)
				return
			}

			tokenAccess, err := tokenReader.ReadTokenAccess(ctx, tokenCookie.Value)
			if err != nil {
				switch err.(type) {
				case *client.ErrUnauthorized:
					handler.HandleError(ctx, w, fmt.Errorf("error checking token access: unauthorized: %v", err), http.StatusUnauthorized, false)
				default:
					handler.HandleError(ctx, w, fmt.Errorf("error checking token access: %v", err), http.StatusInternalServerError, true)
				}
				return
			}

			if tokenAccess < expectedAccess {
				handler.HandleError(ctx, w, errors.New("insufficient access level"), http.StatusForbidden, false)
				return
			}

			// Token is valid, access level is sufficient, proceed to the handler.
			next.ServeHTTP(w, r)
		})
	}
}
