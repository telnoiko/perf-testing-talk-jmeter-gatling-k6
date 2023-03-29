package rest

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
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

var (
	limit = 10
	skip  = 0
)

func (t *task) create() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(UserContextKey).(*store.User)
		task := store.Task{}
		c.Bind(&task)
		task.Owner = user.ID
		log.Printf("creating task %v", task)

		err := t.store.Task.Create(&task)
		if err != nil {
			log.Printf("creating task failed %v", err)
			return c.String(http.StatusInternalServerError, "creating task failed")
		}

		response := SanitizedTask{ID: task.ID, Description: task.Description, Completed: task.Completed}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
}

func (t *task) getAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(UserContextKey).(*store.User)
		log.Printf("getting all tasks for user %d", user.ID)
		completed := c.QueryParams().Get("completed")
		limit := c.QueryParams().Get("limit")
		if limit == "" {
			limit = "10"
		}
		skip := c.QueryParams().Get("skip")
		if skip == "" {
			skip = "0"
		}

		intLimit, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}
		intSkip, err := strconv.Atoi(skip)
		if err != nil {
			return err
		}

		tasks, err := t.store.Task.GetAll(user.ID, completed, intSkip, intLimit)
		if err != nil {
			log.Printf("getting all tasks failed %v", err)
			return c.String(http.StatusInternalServerError, "getting all tasks failed")
		}

		response := make([]SanitizedTask, len(tasks))
		for i, task := range tasks {
			response[i] = SanitizedTask{ID: task.ID, Description: task.Description, Completed: task.Completed}
		}

		return c.JSONPretty(http.StatusOK, response, "  ")
	}
}

func (t *task) getById() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(UserContextKey).(*store.User)
		id := c.Param("id")
		log.Printf("getting task %s", id)

		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("get task by id: couldn't parse id %v", err)
			return c.String(http.StatusBadRequest, "get task by id: couldn't parse id")
		}

		task, err := t.store.Task.GetById(taskId, user.ID)
		if err != nil {
			log.Printf("getting task failed %v", err)
			return c.String(http.StatusBadRequest, "getting task failed")
		}

		response := SanitizedTask{ID: task.ID, Description: task.Description, Completed: task.Completed}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
}

func (t *task) delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(UserContextKey).(*store.User)
		id := c.Param("id")
		log.Printf("deleting task %s", id)

		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("update task: couldn't parse id:  %v", err)
			return c.String(http.StatusBadRequest, "update task: couldn't parse id")
		}

		err = t.store.Task.Delete(taskId, user.ID)
		if err != nil {
			log.Printf("deleting task failed %v", err)
			return c.String(http.StatusInternalServerError, "creating task failed")
		}

		return c.NoContent(http.StatusOK)
	}
}

func (t *task) update() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(UserContextKey).(*store.User)
		id := c.Param("id")
		log.Printf("updating task %s", id)

		task := store.Task{}
		c.Bind(&task)

		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("update task: couldn't parse id:  %v", err)
			return c.String(http.StatusBadRequest, "update task: couldn't parse id")
		}

		updated, err := t.store.Task.Update(taskId, user.ID, &task)
		if err != nil {
			log.Printf("updating task failed %v", err)
			return c.String(http.StatusBadRequest, "updating task failed")
		}

		response := SanitizedTask{ID: updated.ID, Description: updated.Description, Completed: updated.Completed}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
}
