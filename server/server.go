package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-server/client/auth"
	"github.com/goodleby/golang-server/client/database"
	"github.com/goodleby/golang-server/config"
	"github.com/goodleby/golang-server/server/handler"
	"github.com/goodleby/golang-server/server/middleware"
)

const v1API string = "/api/v1"

// Server is a server.
type Server struct {
	Config *config.Config
	Router *chi.Mux
	HTTP   *http.Server
	DB     *database.Client
}

// New creates a new server.
func New(ctx context.Context, config *config.Config) (*Server, error) {
	var s Server

	s.Config = config
	s.Router = chi.NewRouter()
	s.HTTP = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Config.Port),
		Handler: s.Router,
	}

	dbClient, err := database.New(ctx, config)
	if err != nil {
		return &s, fmt.Errorf("error creating database client: %v", err)
	}
	s.DB = dbClient

	s.setupRoutes()

	return &s, nil
}

// Start starts the server.
func (s *Server) Start(ctx context.Context) error {
	errc := make(chan error, 1)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go s.serveHTTP(errc)

	select {
	case sig := <-stop:
		log.Printf("Shutdown signal received: %v", sig)
	case err := <-errc:
		log.Printf("Fatal server error: %v", err)
	case <-ctx.Done():
		log.Printf("Server context cancelled")
	}

	if err := s.stop(ctx); err != nil {
		return fmt.Errorf("error stopping server: %v", err)
	}

	return nil
}

// stop stops the server.
func (s *Server) stop(ctx context.Context) error {
	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("error closing database connection: %v", err)
	}

	if err := s.HTTP.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down http server: %v", err)
	}

	log.Print("Server has been gracefully shut down")

	return nil
}

// serverHTTP starts HTTP server listening on its port.
func (s *Server) serveHTTP(errc chan<- error) {
	log.Printf("Server is listening on port: %d", s.Config.Port)
	if err := s.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		errc <- fmt.Errorf("error listening and serving: %v", err)
	}
}

// setupRoutes sets up the routes for the server.
func (s *Server) setupRoutes() {
	s.Router.Get("/_healthz", handler.Healthz)

	s.Router.Route(v1API, func(r chi.Router) {
		r.Get("/articles", handler.GetAllArticles(s.DB))
		r.Post("/articles", handler.AddArticle(s.DB))

		r.Get("/articles/{id}", handler.GetArticle(s.DB))
		r.Delete("/articles/{id}", handler.RemoveArticle(s.DB))
		r.Patch("/articles/{id}", handler.UpdateArticle(s.DB))
	})
}
