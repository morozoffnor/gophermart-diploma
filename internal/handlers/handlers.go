package handlers

import (
	"context"
	"github.com/morozoffnor/gophermart-diploma/internal/accrual"
	"github.com/morozoffnor/gophermart-diploma/internal/auth"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
	"github.com/morozoffnor/gophermart-diploma/internal/storage"
)

type Handlers struct {
	cfg    *config.Config
	auth   *auth.Auth
	db     *storage.DB
	worker *accrual.Worker
}

type WithdrawRequest struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}

func New(cfg *config.Config, auth *auth.Auth, db *storage.DB, w *accrual.Worker) *Handlers {
	h := &Handlers{
		cfg:    cfg,
		auth:   auth,
		db:     db,
		worker: w,
	}
	go w.Start(context.Background())
	return h
}
