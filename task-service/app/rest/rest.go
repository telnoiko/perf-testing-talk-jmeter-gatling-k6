package rest

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"task-service/app/store"
)

type API struct {
	echo       *echo.Echo
	authorizer *Authorizer
	store      *store.Store

	userRest *user
	taskRest *task
}

func New() *API {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Recover())

	store, err := store.New(e.Logger)
	if err != nil {
		panic(err)
	}

	token := &Token{}
	auth := NewAuth(store, token)

	u := user{
		store:  store,
		token:  token,
		logger: e.Logger,
	}

	t := task{
		store:  store,
		logger: e.Logger,
	}

	api := API{e, auth, store, &u, &t}

	return &api
}

func (s *API) Start(port string) {
	s.echo.GET("/", s.check())
	s.echo.POST("/users", s.userRest.create())
	s.echo.POST("/users/login", s.userRest.login())
	s.echo.POST("/users/logout", s.userRest.logout(), s.authorizer.Authorize)

	//e.GET("/tasks", rest.check(), s.authorizer.Authorize)
	//e.GET("/tasks/:id", rest.check(), s.authorizer.Authorize)
	s.echo.POST("/tasks", s.taskRest.create(), s.authorizer.Authorize)
	//e.PUT("/tasks/:id", rest.check(), s.authorizer.Authorize)
	//e.DELETE("/tasks/:id", rest.check(), s.authorizer.Authorize)

	if err := s.echo.Start(port); err != nil && err != http.ErrServerClosed {
		s.echo.Logger.Fatal("shutting down the server", err)
	}
}

func (s *API) check() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}
}

func (s *API) Stop(ctx context.Context) {
	s.store.Stop()
	if err := s.echo.Shutdown(ctx); err != nil {
		s.echo.Logger.Fatal(err)
	}
}
