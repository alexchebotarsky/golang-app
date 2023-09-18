package main

import (
	"context"
	"log"

	"github.com/goodleby/golang-server/app"
	"github.com/goodleby/golang-server/config"
	"github.com/goodleby/golang-server/metrics"
	"github.com/goodleby/golang-server/tracing"
)

func main() {
	ctx := context.Background()

	config, err := config.Load(ctx, ".env")
	if err != nil {
		log.Fatalf("Error loading env config: %v", err)
	}

	if err := tracing.Init(config); err != nil {
		log.Fatalf("Error initializing tracing: %v", err)
	}

	if err := metrics.Init(); err != nil {
		log.Fatalf("Error initializing metrics: %v", err)
	}

	s, err := app.New(ctx, config)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	s.Start(ctx)
}
