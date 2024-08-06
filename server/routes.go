package server

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/goodleby/golang-app/client/auth"
	"github.com/goodleby/golang-app/server/handler"
	"github.com/goodleby/golang-app/server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) setupRoutes(allowedOrigins []string) {
	s.Router.Get("/_healthz", handler.Health)
	s.Router.Handle("/metrics", promhttp.Handler())

	s.Router.Route(v1API, func(r chi.Router) {
		r.Use(middleware.Trace, middleware.Metrics)

		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   allowedOrigins,
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
			AllowCredentials: true,
		}))

		r.Get("/example", handler.GetExampleData(s.Clients.Example))
		r.Post("/pubsub/articles", handler.AddArticlePubSub(s.Clients.PubSub))

		// Auth routes
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", handler.AuthLogin(s.Clients.Auth))
			r.Post("/auth/refresh", handler.AuthRefresh(s.Clients.Auth))
			r.Post("/auth/logout", handler.AuthLogout)
		})

		// View articles
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(s.Clients.Auth, auth.ViewerAccess))

			r.Get("/articles", handler.GetAllArticles(s.Clients.DB))
			r.Get("/articles/{id}", handler.GetArticle(s.Clients.DB))
		})

		// Edit articles
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(s.Clients.Auth, auth.EditorAccess))

			r.Post("/articles", handler.AddArticle(s.Clients.DB))
			r.Delete("/articles/{id}", handler.DeleteArticle(s.Clients.DB))
			r.Put("/articles/{id}", handler.UpdateArticle(s.Clients.DB))
		})
	})
}

const v1API string = "/api/v1"
