package main

import (
	"context"
	"log"

	"github.com/goodleby/golang-server/app"
	"github.com/goodleby/golang-server/env"
	"github.com/goodleby/golang-server/metrics"
	"github.com/goodleby/golang-server/tracing"
)

func main() {
	ctx := context.Background()

	env, err := env.LoadConfig(ctx, ".env")
	if err != nil {
		log.Fatalf("Error loading env config: %v", err)
	}

	if err := tracing.Init(env); err != nil {
		log.Fatalf("Error initializing tracing: %v", err)
	}

	if err := metrics.Init(); err != nil {
		log.Fatalf("Error initializing metrics: %v", err)
	}

	s, err := app.New(ctx, env)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	s.Start(ctx)
}
