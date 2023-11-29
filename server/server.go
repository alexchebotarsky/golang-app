package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-app/client/auth"
	"github.com/goodleby/golang-app/server/handler"
	"github.com/goodleby/golang-app/server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
	middleware.TokenParser
}

type PubSubClient interface {
	handler.AddArticlePublisher
}

type ExampleClient interface {
	handler.ExampleDataFetcher
}

type Server struct {
	Port    uint16
	Router  chi.Router
	HTTP    *http.Server
	DB      DBClient
	Auth    AuthClient
	PubSub  PubSubClient
	Example ExampleClient
}

func New(ctx context.Context, port uint16, db DBClient, auth AuthClient, pubsub PubSubClient, example ExampleClient) (*Server, error) {
	s := &Server{}

	s.Port = port
	s.Router = chi.NewRouter()
	s.HTTP = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Port),
		Handler:      s.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	s.DB = db
	s.Auth = auth
	s.PubSub = pubsub
	s.Example = example

	s.setupRoutes()

	return s, nil
}

func (s *Server) Start(ctx context.Context) {
	slog.Info(fmt.Sprintf("Server is listening on port: %d", s.Port))
	if err := s.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error(fmt.Sprintf("Error listening and serving: %v", err))
	}
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.HTTP.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down http server: %v", err)
	}

	slog.Debug("Server has stopped gracefully")

	return nil
}

func (s *Server) setupRoutes() {
	s.Router.Get("/_healthz", handler.Health)
	s.Router.Handle("/metrics", promhttp.Handler())

	s.Router.Route(v1API, func(r chi.Router) {
		r.Use(middleware.Trace, middleware.Metrics)

		r.Get("/example", handler.GetExampleData(s.Example))
		r.Post("/pubsub/articles", handler.AddArticlePubSub(s.PubSub))

		// Auth routes
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", handler.AuthLogin(s.Auth))
			r.Post("/auth/refresh", handler.AuthRefresh(s.Auth))
			r.Post("/auth/logout", handler.AuthLogout)
		})

		// View articles
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(s.Auth, auth.ViewerAccess))

			r.Get("/articles", handler.GetAllArticles(s.DB))
			r.Get("/articles/{id}", handler.GetArticle(s.DB))
		})

		// Edit articles
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(s.Auth, auth.EditorAccess))

			r.Post("/articles", handler.AddArticle(s.DB))
			r.Delete("/articles/{id}", handler.DeleteArticle(s.DB))
			r.Put("/articles/{id}", handler.UpdateArticle(s.DB))
		})
	})
}

const v1API string = "/api/v1"
