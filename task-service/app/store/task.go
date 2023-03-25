package store

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	UserID      int    `json:"user_id"`
}

type TaskStore struct {
	pool   *pgxpool.Pool
	logger echo.Logger
}

func (s *TaskStore) Create(task *Task) error {

	err := s.pool.QueryRow(context.Background(), `INSERT INTO tasks (description, completed, user_id) VALUES ($1, $2, $3) RETURNING id`, task.Description, task.Completed, task.UserID).Scan(&task.ID)
	if err != nil {
		s.logger.Printf("task create failed: %v\n", err)
		return err
	}

	return nil
}
