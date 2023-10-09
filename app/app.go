package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/goodleby/golang-server/client/auth"
	"github.com/goodleby/golang-server/client/database"
	"github.com/goodleby/golang-server/client/example"
	"github.com/goodleby/golang-server/env"
	"github.com/goodleby/golang-server/server"
)

type Service interface {
	Run(context.Context)
}

type App struct {
	Services []Service
	Env      *env.Config
	DB       *database.Client
	Auth     *auth.Client
	Example  *example.Client
}

func New(ctx context.Context, env *env.Config) (*App, error) {
	app := &App{
		Env: env,
	}

	dbClient, err := database.New(ctx, database.Credentials{
		User:     env.DatabaseUser,
		Password: env.DatabasePassword,
		Host:     env.DatabaseHost,
		Port:     env.DatabasePort,
		Name:     env.DatabaseName,
		Options:  env.DatabaseOptions,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating database client: %v", err)
	}
	app.DB = dbClient

	authClient, err := auth.New(ctx, env.AuthSecret, auth.Keys{
		Admin:  env.AuthAdminKey,
		Editor: env.AuthEditorKey,
		Viewer: env.AuthViewerKey,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating auth client: %v", err)
	}
	app.Auth = authClient

	exampleClient, err := example.New(env.ExampleEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating example client: %v", err)
	}
	app.Example = exampleClient

	server, err := server.New(ctx, env.Port, app.DB, app.Auth, app.Example)
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
