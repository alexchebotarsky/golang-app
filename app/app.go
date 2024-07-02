package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/goodleby/golang-app/client"
	"github.com/goodleby/golang-app/client/auth"
	"github.com/goodleby/golang-app/client/database"
	"github.com/goodleby/golang-app/client/example"
	"github.com/goodleby/golang-app/client/pubsub"
	"github.com/goodleby/golang-app/env"
	"github.com/goodleby/golang-app/processor"
	"github.com/goodleby/golang-app/server"
)

type App struct {
	Services []Service
	Clients  *Clients
}

type Service interface {
	Start(context.Context)
	Stop(context.Context) error
}

type Clients struct {
	DB      *database.Client
	Auth    *auth.Client
	PubSub  *pubsub.Client
	Example *example.Client
}

func (c *Clients) Close() error {
	var errors []error

	slog.Debug("Clients are closing...")

	err := c.DB.Close()
	if err != nil {
		errors = append(errors, fmt.Errorf("error closing database client: %v", err))
	}

	if len(errors) > 0 {
		return &client.ErrMultiple{Errs: errors}
	}

	slog.Debug("Clients closing complete")

	return nil
}

func New(ctx context.Context, env *env.Config) (*App, error) {
	var app App
	var err error

	app.Clients, err = setupClients(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("error setting up clients: %v", err)
	}

	app.Services, err = setupServices(ctx, env, app.Clients)
	if err != nil {
		return nil, fmt.Errorf("error setting up services: %v", err)
	}

	return &app, nil
}

func (app *App) Launch(ctx context.Context) {
	for _, service := range app.Services {
		go service.Start(ctx)
	}

	<-ctx.Done()

	slog.Debug("App is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, service := range app.Services {
		err := service.Stop(ctx)
		if err != nil {
			slog.Error(fmt.Sprintf("Error stopping a service: %v", err))
		}
	}

	err := app.Clients.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("Error closing app clients: %v", err))
	}

	slog.Debug("App shutdown complete")
}

func setupClients(ctx context.Context, env *env.Config) (*Clients, error) {
	var c Clients
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

	c.Auth = auth.New(ctx, env.AuthSecret, env.AuthTokenTTL, auth.Keys{
		Admin:  env.AuthAdminKey,
		Editor: env.AuthEditorKey,
		Viewer: env.AuthViewerKey,
	})

	c.PubSub, err = pubsub.New(ctx, env.PubSubProjectID, env.Environment)
	if err != nil {
		return nil, fmt.Errorf("error creating example client: %v", err)
	}

	c.Example = example.New(env.ExampleEndpoint)

	return &c, nil
}

func setupServices(ctx context.Context, env *env.Config, clients *Clients) ([]Service, error) {
	var services []Service

	server, err := server.New(ctx, env.Port, env.AllowedOrigins, server.Clients{
		DB:      clients.DB,
		Auth:    clients.Auth,
		PubSub:  clients.PubSub,
		Example: clients.Example,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating new server: %v", err)
	}
	services = append(services, server)

	processor, err := processor.New(ctx, processor.Clients{
		PubSub: clients.PubSub,
		DB:     clients.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating new processor: %v", err)
	}
	services = append(services, processor)

	return services, nil
}
