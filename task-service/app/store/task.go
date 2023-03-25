package store

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"strconv"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Owner       int    `json:"owner"`
}

type TaskStore struct {
	pool   *pgxpool.Pool
	logger echo.Logger
}

func (s *TaskStore) Create(task *Task) error {

	row := s.pool.QueryRow(context.Background(), `INSERT INTO tasks (description, completed, owner) VALUES ($1, $2, $3) RETURNING id`, task.Description, task.Completed, task.Owner)
	err := row.Scan(&task.ID)
	if err != nil {
		s.logger.Printf("task create failed: %v\n", err)
		return err
	}

	return nil
}

func (s *TaskStore) GetAll(userID int, completed string, skip int, limit int) ([]*Task, error) {
	var rows pgx.Rows
	var err error

	if completed != "" {
		parseBool, err := strconv.ParseBool(completed)
		if err != nil {
			return nil, err
		}
		rows, err = s.pool.Query(context.Background(), `SELECT id, description, completed FROM tasks WHERE owner = $1 AND completed = $2 ORDER BY id LIMIT $3 OFFSET $4`, userID, parseBool, limit, skip)
	} else {
		rows, err = s.pool.Query(context.Background(), `SELECT id, description, completed FROM tasks WHERE owner = $1 ORDER BY id LIMIT $2 OFFSET $3`, userID, limit, skip)
	}
	if err != nil {
		s.logger.Printf("get all tasks failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Description, &task.Completed)
		if err != nil {
			s.logger.Printf("get all tasks failed: %v\n", err)
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *TaskStore) GetById(taskID int, userID int) (*Task, error) {
	row := s.pool.QueryRow(context.Background(), `SELECT id, description, completed FROM tasks WHERE id = $1 AND owner = $2`, taskID, userID)
	task := &Task{}
	err := row.Scan(&task.ID, &task.Description, &task.Completed)
	if err != nil {
		s.logger.Printf("get task by id failed: %v\n", err)
		return nil, err
	}
	return task, nil
}

func (s *TaskStore) Update(taskID int, userID int, task *Task) (*Task, error) {
	row := s.pool.QueryRow(context.Background(), `UPDATE tasks SET description = $1, completed = $2 WHERE id = $3 AND owner = $4 RETURNING id, description, completed`, task.Description, task.Completed, taskID, userID)
	updatedTask := &Task{}
	err := row.Scan(&updatedTask.ID, &updatedTask.Description, &updatedTask.Completed)
	if err != nil {
		s.logger.Printf("task update failed: %v\n", err)
		return nil, err
	}
	return updatedTask, nil
}

func (s *TaskStore) Delete(taskID int, userID int) error {
	_, err := s.pool.Exec(context.Background(), `DELETE FROM tasks WHERE id = $1 AND owner = $2`, taskID, userID)
	if err != nil {
		s.logger.Printf("task delete failed: %v\n", err)
		return err
	}
	return nil
}
