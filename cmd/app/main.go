package main

import (
	"context"
	"log"

	"github.com/goodleby/golang-app/app"
	"github.com/goodleby/golang-app/env"
	"github.com/goodleby/golang-app/metrics"
	"github.com/goodleby/golang-app/tracing"
)

func main() {
	ctx := context.Background()

	env, err := env.LoadConfig(ctx, ".env")
	if err != nil {
		log.Fatalf("Error loading env config: %v", err)
	}

	if err := tracing.Init(env.ServiceName); err != nil {
		log.Fatalf("Error initializing tracing: %v", err)
	}

	if err := metrics.Init(); err != nil {
		log.Fatalf("Error initializing metrics: %v", err)
	}

	app, err := app.New(ctx, env)
	if err != nil {
		log.Fatalf("Error creating app: %v", err)
	}

	app.Launch(ctx)
}
