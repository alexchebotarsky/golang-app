package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-server/client/auth"
	"github.com/goodleby/golang-server/server/handler"
)

type TokenParser interface {
	ParseToken(ctx context.Context, token string) (*auth.Claims, error)
}

func Auth(tokenParser TokenParser, expectedAccessLevel int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tokenCookie, err := r.Cookie("token")
			if err != nil {
				handler.HandleError(ctx, w, fmt.Errorf("error reading auth token cookie: %v", err), http.StatusUnauthorized, true)
				return
			}

			claims, err := tokenParser.ParseToken(ctx, tokenCookie.Value)
			if err != nil {
				handler.HandleError(ctx, w, fmt.Errorf("error validating auth token: %v", err), http.StatusUnauthorized, true)
				return
			}

			if expectedAccessLevel > claims.AccessLevel {
				handler.HandleError(ctx, w, errors.New("insufficient access level"), http.StatusForbidden, true)
				return
			}

			// Token is valid, access level is sufficient, proceed to the handler.
			next.ServeHTTP(w, r)
		})
	}
}
