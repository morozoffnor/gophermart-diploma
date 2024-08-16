package main

import (
	"context"
	"fmt"
	"github.com/morozoffnor/gophermart-diploma/internal/accrual"
	"github.com/morozoffnor/gophermart-diploma/internal/auth"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
	"github.com/morozoffnor/gophermart-diploma/internal/handlers"
	"github.com/morozoffnor/gophermart-diploma/internal/middlewares"
	"github.com/morozoffnor/gophermart-diploma/internal/server"
	"github.com/morozoffnor/gophermart-diploma/internal/storage"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.New()
	a := auth.New(cfg)
	db := storage.New(cfg, ctx)
	c := accrual.NewClient(cfg)
	w := accrual.NewWorker(cfg, db, c)
	h := handlers.New(cfg, a, db, w)
	m := middlewares.New(a, db)
	r := server.NewRouter(h, m)
	s := server.NewSever(cfg, r)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	wg, gCtx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return s.ListenAndServe()
	})
	wg.Go(func() error {
		<-gCtx.Done()
		return s.Shutdown(context.Background())
	})

	// если зарегать заказ и, не дождавшись обработки от сервиса начислений, выключить сервис,
	// то заказ так и останется висеть необработанным в бд. в тз не было такого требования,
	// но мне показалось логичным при старте сервиса проверять есть ли необработанные заказы
	go w.ProcessStaleOrders(ctx)

	if err := wg.Wait(); err != nil {
		fmt.Printf("exit reason: %s", err)
	}
}
