package store

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Store struct {
	pool   *pgxpool.Pool
	logger echo.Logger

	User *UserStore
	Task *TaskStore
}

func New(logger echo.Logger) (*Store, error) {
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:password@db:5432/tasks?pool_max_conns=10")
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		return nil, err
	}

	user := &UserStore{
		pool:   pool,
		logger: logger,
	}

	task := &TaskStore{
		pool:   pool,
		logger: logger,
	}

	return &Store{pool, logger, user, task}, nil
}

func (s *Store) Stop() {
	s.pool.Close()
}
