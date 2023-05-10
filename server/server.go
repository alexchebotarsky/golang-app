package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/pure-go-server/client/database"
	"github.com/goodleby/pure-go-server/config"
	"github.com/goodleby/pure-go-server/server/handler"
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
func (s *Server) Start() error {
	log.Printf("Server is listening on port :%d", s.Config.Port)
	if err := s.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("unexpected server error: %v", err)
	}

	return nil
}

// setupRoutes sets up the routes for the server.
func (s *Server) setupRoutes() {
	s.Router.Get("/_healthz", handler.Healthz)

	s.Router.Route(v1API, func(r chi.Router) {
		r.Get("/articles", handler.GetAllArticles(s.DB))
		r.Post("/articles", handler.CreateArticle(s.DB))

		r.Get("/articles/{id}", handler.GetArticle(s.DB))
		r.Delete("/articles/{id}", handler.RemoveArticle(s.DB))
		r.Patch("/articles/{id}", handler.UpdateArticle(s.DB))
	})
}
