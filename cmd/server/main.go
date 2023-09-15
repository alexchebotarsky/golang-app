package main

import (
	"context"
	"log"
	"os"

	"github.com/goodleby/golang-server/config"
	"github.com/goodleby/golang-server/metrics"
	"github.com/goodleby/golang-server/server"
	"github.com/goodleby/golang-server/tracing"
)

func main() {
	ctx := context.Background()

	config, err := config.Load(ctx, ".env")
	if err != nil {
		log.Printf("Error loading config: %v", err)
		os.Exit(1)
	}

	if err := tracing.Init(config); err != nil {
		log.Printf("Error initializing tracing: %v", err)
		os.Exit(1)
	}

	if err := metrics.Init(); err != nil {
		log.Printf("Error initializing metrics: %v", err)
		os.Exit(1)
	}

	s, err := server.New(ctx, config)
	if err != nil {
		log.Printf("Error creating server: %v", err)
		os.Exit(1)
	}

	s.Start(ctx)
}
