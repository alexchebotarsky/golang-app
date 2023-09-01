package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-server/client/auth"
	"github.com/goodleby/golang-server/client/database"
	"github.com/goodleby/golang-server/client/example"
	"github.com/goodleby/golang-server/config"
	"github.com/goodleby/golang-server/metrics"
	"github.com/goodleby/golang-server/server/handler"
	"github.com/goodleby/golang-server/server/middleware"
)

const v1API string = "/api/v1"

// Server is a server.
type Server struct {
	Config  *config.Config
	Router  *chi.Mux
	HTTP    *http.Server
	DB      *database.Client
	Auth    *auth.Client
	Example *example.Client
}

// New creates a new server.
func New(ctx context.Context, config *config.Config) (*Server, error) {
	var s Server

	s.Config = config
	s.Router = chi.NewRouter()
	s.HTTP = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.Port),
		Handler:      s.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	dbClient, err := database.New(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating database client: %v", err)
	}
	s.DB = dbClient

	authClient, err := auth.New(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating database client: %v", err)
	}
	s.Auth = authClient

	exampleClient, err := example.New(config)
	if err != nil {
		return nil, fmt.Errorf("error creating database client: %v", err)
	}
	s.Example = exampleClient

	if err := metrics.Init(); err != nil {
		return nil, fmt.Errorf("error initializing metrics: %v", err)
	}

	s.setupRoutes()

	return &s, nil
}

// Start starts the server.
func (s *Server) Start(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	go s.serveHTTP()

	<-ctx.Done()

	ctx, cancelTimeout := context.WithTimeout(ctx, 10*time.Second)
	defer cancelTimeout()

	s.Stop(ctx)
}

// Stop stops the server.
func (s *Server) Stop(ctx context.Context) {
	if err := s.HTTP.Shutdown(ctx); err != nil {
		log.Printf("error shutting down http server: %v", err)
	}

	if err := s.DB.Close(); err != nil {
		log.Printf("error closing database connection: %v", err)
	}

	log.Print("Server has been gracefully stopped")
}

// serverHTTP starts HTTP server listening on its port.
func (s *Server) serveHTTP() {
	log.Printf("Server is listening on port: %d", s.Config.Port)
	if err := s.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Error listening and serving: %v", err)
	}
}

// setupRoutes sets up the routes for the server.
func (s *Server) setupRoutes() {
	s.Router.Get("/_healthz", handler.Healthz)

	s.Router.Route(v1API, func(r chi.Router) {
		r.Use(middleware.Trace)
		r.Use(middleware.Metrics)

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
			r.Delete("/articles/{id}", handler.RemoveArticle(s.DB))
			r.Patch("/articles/{id}", handler.UpdateArticle(s.DB))
		})
	})
}
