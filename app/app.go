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

type Clients struct {
	DB      *database.Client
	Auth    *auth.Client
	Example *example.Client
}

type App struct {
	Services []Service
}

func New(ctx context.Context, env *env.Config) (*App, error) {
	app := &App{}

	clients, err := setupClients(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("error setting up clients: %v", err)
	}

	app.Services, err = setupServices(ctx, env, clients)
	if err != nil {
		return nil, fmt.Errorf("error setting up services: %v", err)
	}

	return app, nil
}

func (app *App) Launch(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	for _, service := range app.Services {
		go service.Run(ctx)
	}
}

func setupClients(ctx context.Context, env *env.Config) (*Clients, error) {
	c := &Clients{}
	var err error

	c.DB, err = database.New(ctx, database.Credentials{
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

	c.Auth, err = auth.New(ctx, env.AuthSecret, auth.Keys{
		Admin:  env.AuthAdminKey,
		Editor: env.AuthEditorKey,
		Viewer: env.AuthViewerKey,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating auth client: %v", err)
	}

	c.Example, err = example.New(env.ExampleEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating example client: %v", err)
	}

	return c, nil
}

func setupServices(ctx context.Context, env *env.Config, clients *Clients) ([]Service, error) {
	services := []Service{}

	server, err := server.New(ctx, env.Port, clients.DB, clients.Auth, clients.Example)
	if err != nil {
		return nil, fmt.Errorf("error creating new server: %v", err)
	}
	services = append(services, server)

	return services, nil
}
