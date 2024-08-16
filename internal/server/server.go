package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
	"github.com/morozoffnor/gophermart-diploma/internal/handlers"
	"github.com/morozoffnor/gophermart-diploma/internal/middlewares"
	"log"
	"net/http"
)

func NewRouter(h *handlers.Handlers, m *middlewares.Middlewares) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Group(func(r chi.Router) {
		// эндпоинты без авторизационной миддлвари
		r.Post("/api/user/register", h.RegisterUser)
		r.Post("/api/user/login", h.LoginUser)
	})
	r.Group(func(r chi.Router) {
		r.Use(m.Auth())
		// эндпоинты с авторизацией
		r.Get("/api/user/orders", h.GetOrders)
		r.Post("/api/user/orders", h.UploadOrder)
		r.Post("/api/user/balance/withdraw", h.Withdraw)
		r.Get("/api/user/balance", h.GetBalance)
		r.Get("/api/user/withdrawals", h.GetWithdrawals)
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
