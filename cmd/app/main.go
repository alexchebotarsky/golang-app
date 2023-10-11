package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/goodleby/golang-app/app"
	"github.com/goodleby/golang-app/env"
	"github.com/goodleby/golang-app/logger"
	"github.com/goodleby/golang-app/metrics"
	"github.com/goodleby/golang-app/tracing"
)

func main() {
	ctx := context.Background()

	env, err := env.LoadConfig(ctx, ".env")
	if err != nil {
		slog.Error(fmt.Sprintf("Error loading env config: %v", err))
		os.Exit(1)
	}

	logger.Init()

	if err := tracing.Init(env.ServiceName); err != nil {
		slog.Error(fmt.Sprintf("Error initializing tracing: %v", err))
		os.Exit(1)
	}

	if err := metrics.Init(); err != nil {
		slog.Error(fmt.Sprintf("Error initializing metrics: %v", err))
		os.Exit(1)
	}

	app, err := app.New(ctx, env)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating app: %v", err))
		os.Exit(1)
	}

	app.Launch(ctx)
}
