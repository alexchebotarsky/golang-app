package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/goodleby/golang-app/app"
	"github.com/goodleby/golang-app/env"
	"github.com/goodleby/golang-app/logger"
	"github.com/goodleby/golang-app/metrics"
	"github.com/goodleby/golang-app/tracing"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	env, err := env.LoadConfig(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("Error loading env config: %v", err))
		os.Exit(1)
	}

	logger.Init(env.LogLevel, env.LogFormat)

	err = tracing.Init(ctx, env.ServiceName, env.Environment, env.TracingSampleRate)
	if err != nil {
		slog.Error(fmt.Sprintf("Error initializing tracing: %v", err))
	}

	err = metrics.Init()
	if err != nil {
		slog.Error(fmt.Sprintf("Error initializing metrics: %v", err))
	}

	app, err := app.New(ctx, env)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating app: %v", err))
		os.Exit(1)
	}

	app.Launch(ctx)
}
