package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
	"github.com/morozoffnor/gophermart-diploma/internal/handlers"
	"log"
	"net/http"
)

func NewRouter(h *handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Group(func(r chi.Router) {
		// No-auth endpoints
		r.Post("/api/user/register", h.RegisterUser)
		r.Post("/api/user/login", h.LoginUser)
	})
	r.Group(func(r chi.Router) {
		// Auth endpoints
	})
	return r
}

func NewSever(cfg *config.Config, r *chi.Mux) *http.Server {
	s := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}
	defer log.Printf("\nServer started on addr: %s\nDatabase URI: %s\nAccrual Addr: %s",
		cfg.Addr, cfg.DatabaseURI, cfg.AccrualSystemAddr)
	return s
}
