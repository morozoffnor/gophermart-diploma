package handlers

import (
	"github.com/morozoffnor/gophermart-diploma/internal/auth"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
	"github.com/morozoffnor/gophermart-diploma/internal/storage"
)

type Handlers struct {
	cfg  *config.Config
	auth *auth.Auth
	db   *storage.DB
}

func New(cfg *config.Config, auth *auth.Auth, db *storage.DB) *Handlers {
	h := &Handlers{
		cfg:  cfg,
		auth: auth,
		db:   db,
	}
	return h
}
