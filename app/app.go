package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/goodleby/golang-server/env"
	"github.com/goodleby/golang-server/server"
)

type Service interface {
	Run(context.Context)
}

type App struct {
	Services []Service
	Env      *env.Config
}

func New(ctx context.Context, env *env.Config) (*App, error) {
	app := &App{}

	app.Env = env

	server, err := server.New(ctx, app.Env)
	if err != nil {
		return nil, fmt.Errorf("error creating new server: %v", err)
	}
	app.Services = append(app.Services, server)

	return app, nil
}

func (app *App) Launch(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	for _, service := range app.Services {
		go service.Run(ctx)
	}
}
