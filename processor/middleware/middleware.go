package middleware

import "github.com/goodleby/golang-app/processor/event"

type Middleware func(eventID string, next event.Handler) event.Handler
