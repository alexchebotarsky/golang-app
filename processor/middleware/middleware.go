package middleware

import "github.com/goodleby/golang-server/processor/event"

type Middleware func(eventID string, next event.Handler) event.Handler
