package rest

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"task-service/app/store"
)

type task struct {
	store  *store.Store
	logger echo.Logger
}

type SanitizedTask struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (t *task) create() func(c echo.Context) error {
	return func(c echo.Context) error {
		user := c.Get(UserContextKey).(*store.User)
		task := store.Task{}
		c.Bind(&task)
		task.UserID = user.ID
		t.logger.Infof("creating task %v", task)

		err := t.store.Task.Create(&task)
		if err != nil {
			return err
		}

		response := SanitizedTask{ID: task.ID, Description: task.Description, Completed: task.Completed}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
}
