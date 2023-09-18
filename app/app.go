package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/goodleby/golang-server/config"
	"github.com/goodleby/golang-server/server"
)

type Service interface {
	Start(context.Context)
	Stop(context.Context)
}

type App struct {
	Services []Service
}

func New(ctx context.Context, config *config.Config) (*App, error) {
	var app App

	server, err := server.New(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating new server: %v", err)
	}
	app.Services = append(app.Services, server)

	return &app, nil
}

func (app *App) Start(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	for _, service := range app.Services {
		go service.Start(ctx)
	}

	<-ctx.Done()

	ctx, cancelTimeout := context.WithTimeout(ctx, 10*time.Second)
	defer cancelTimeout()

	app.Stop(ctx)
}

func (app *App) Stop(ctx context.Context) {
	for _, service := range app.Services {
		service.Stop(ctx)
	}
	log.Print("App has been gracefully stopped")
}
