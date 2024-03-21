package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-app/client"
	"github.com/goodleby/golang-app/server/handler"
)

type TokenChecker interface {
	CheckTokenAccess(ctx context.Context, token string, expectedAccessLevel int) error
}

func Auth(tokenChecker TokenChecker, expectedAccessLevel int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tokenCookie, err := r.Cookie("token")
			if err != nil {
				handler.HandleError(ctx, w, fmt.Errorf("error reading auth token cookie: %v", err), http.StatusUnauthorized, false)
				return
			}

			err = tokenChecker.CheckTokenAccess(ctx, tokenCookie.Value, expectedAccessLevel)
			if err != nil {
				switch err.(type) {
				case *client.ErrUnauthorized:
					handler.HandleError(ctx, w, fmt.Errorf("error checking token access: unauthorized: %v", err), http.StatusUnauthorized, false)
				case *client.ErrForbidden:
					handler.HandleError(ctx, w, fmt.Errorf("error checking token access: forbidden: %v", err), http.StatusForbidden, false)
				default:
					handler.HandleError(ctx, w, fmt.Errorf("error checking token access: %v", err), http.StatusInternalServerError, true)
				}
				return
			}

			// Token is valid, access level is sufficient, proceed to the handler.
			next.ServeHTTP(w, r)
		})
	}
}
