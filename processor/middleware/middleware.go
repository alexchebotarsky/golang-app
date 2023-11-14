package middleware

import "github.com/goodleby/golang-app/processor/event"

type Middleware func(eventName string, next event.Handler) event.Handler
