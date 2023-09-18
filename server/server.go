package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-server/client/auth"
	"github.com/goodleby/golang-server/client/database"
	"github.com/goodleby/golang-server/client/example"
	"github.com/goodleby/golang-server/config"
	"github.com/goodleby/golang-server/server/handler"
	"github.com/goodleby/golang-server/server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Config  *config.Config
	Router  chi.Router
	HTTP    *http.Server
	DB      *database.Client
	Auth    *auth.Client
	Example *example.Client
}

func New(ctx context.Context, config *config.Config) (*Server, error) {
	var s Server
	var err error

	s.Config = config
	s.Router = chi.NewRouter()
	s.HTTP = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.Port),
		Handler:      s.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	s.DB, err = database.New(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating database client: %v", err)
	}

	s.Auth, err = auth.New(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating auth client: %v", err)
	}

	s.Example, err = example.New(config)
	if err != nil {
		return nil, fmt.Errorf("error creating example client: %v", err)
	}

	s.setupRoutes()

	return &s, nil
}

func (s *Server) Start(ctx context.Context) {
	log.Printf("Server is listening on port: %d", s.Config.Port)
	if err := s.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Error listening and serving: %v", err)
	}
}

func (s *Server) Stop(ctx context.Context) {
	if err := s.HTTP.Shutdown(ctx); err != nil {
		log.Printf("error shutting down http server: %v", err)
	}
}

func (s *Server) setupRoutes() {
	s.Router.Get("/_healthz", handler.Health)
	s.Router.Handle("/metrics", promhttp.Handler())

	s.Router.Route(v1API, func(r chi.Router) {
		r.Use(middleware.Trace, middleware.Metrics)

		r.Get("/example", handler.GetExampleData(s.Example))

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
			r.Patch("/articles/{id}", handler.UpdateArticle(s.DB))
		})
	})
}

const v1API string = "/api/v1"
