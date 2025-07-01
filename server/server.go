package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-app/server/handler"
	"github.com/goodleby/golang-app/server/middleware"
)

type Server struct {
	Host    string
	Port    uint16
	Router  chi.Router
	HTTP    *http.Server
	Clients Clients
}

type Clients struct {
	DB      DBClient
	Auth    AuthClient
	PubSub  PubSubClient
	Example ExampleClient
}

type DBClient interface {
	handler.AllArticlesSelector
	handler.ArticleSelector
	handler.ArticleInserter
	handler.ArticleUpdater
	handler.ArticleDeleter
}

type AuthClient interface {
	handler.TokenCreator
	handler.TokenRefresher
	middleware.TokenAccessReader
}

type PubSubClient interface {
	handler.AddArticlePublisher
}

type ExampleClient interface {
	handler.ExampleDataFetcher
}

func New(ctx context.Context, host string, port uint16, allowedOrigins []string, clients Clients) (*Server, error) {
	var s Server

	s.Host = host
	s.Port = port
	s.Router = chi.NewRouter()
	s.HTTP = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.Host, s.Port),
		Handler:      s.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	s.Clients = clients

	s.setupRoutes(allowedOrigins)

	return &s, nil
}

func (s *Server) Start(ctx context.Context, errc chan<- error) {
	slog.Info(fmt.Sprintf("Server is listening at %s:%d", s.Host, s.Port))
	err := s.HTTP.ListenAndServe()
	if err != http.ErrServerClosed {
		errc <- fmt.Errorf("error listening and serving: %v", err)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.HTTP.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down http server: %v", err)
	}

	return nil
}
