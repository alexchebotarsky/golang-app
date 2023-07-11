package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-server/client/auth"
	"github.com/goodleby/golang-server/server/handler"
)

// TokenParser is an interface for parsing and validating a JWT.
type TokenParser interface {
	ParseToken(token string) (*auth.Claims, error)
}

// Auth is a middleware that checks authorization cookie and if access level is not sufficient blocks request.
func Auth(tokenParser TokenParser, expectedAccessLevel int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("token")
			if err != nil {
				handler.HandleError(w, fmt.Errorf("error reading auth token cookie: %v", err), http.StatusUnauthorized, true)
				return
			}

			claims, err := tokenParser.ParseToken(tokenCookie.Value)
			if err != nil {
				handler.HandleError(w, fmt.Errorf("error validating auth token: %v", err), http.StatusUnauthorized, true)
				return
			}

			if expectedAccessLevel < claims.AccessLevel {
				handler.HandleError(w, errors.New("insufficient access level"), http.StatusForbidden, true)
				return
			}

			// Token is valid, access level is sufficient, proceed to the handler.
			next.ServeHTTP(w, r)
		})
	}
}
