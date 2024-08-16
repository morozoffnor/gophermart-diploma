package storage

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
	"time"
)

type DB struct {
	pool *pgxpool.Pool
	cfg  *config.Config
}

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UserID     string
	UploadedAt time.Time `json:"uploaded_ats"`
}

func New(cfg *config.Config, ctx context.Context) *DB {
	db := &DB{
		cfg: cfg,
	}
	conn, err := pgxpool.New(ctx, cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}
	conn.Config().MaxConns = 20
	conn.Config().MinConns = 2
	db.pool = conn
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	doMigrations(cfg)

	return db
}

func doMigrations(cfg *config.Config) {
	m, err := migrate.New("file://internal/storage/migrations", cfg.DatabaseURI)

	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}
