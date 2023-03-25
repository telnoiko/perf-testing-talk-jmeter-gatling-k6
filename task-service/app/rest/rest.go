package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"task-service/app/store"
)

type API struct {
	echo       *echo.Echo
	userRest   *user
	authorizer *Authorizer
	store      *store.Store
}

func New(auth *Authorizer, store *store.Store) *API {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Recover())

	u := user{
		store:      store,
		authorizer: auth,
		logger:     e.Logger,
	}

	api := API{e, &u, auth, store}

	return &api
}

func (s *API) Start() {
	s.echo.GET("/", s.check())
	s.echo.POST("/users", s.userRest.create())
	s.echo.POST("/users/login", s.userRest.login())
	s.echo.POST("/users/logout", s.userRest.logout(), s.authorizer.Authorize)

	//e.GET("/tasks", rest.check(), s.authorizer.Authorize)
	//e.GET("/tasks/:id", rest.check(), s.authorizer.Authorize)
	//e.POST("/tasks", s.c(), s.authorizer.Authorize)
	//e.PUT("/tasks/:id", rest.check(), s.authorizer.Authorize)
	//e.DELETE("/tasks/:id", rest.check(), s.authorizer.Authorize)

	s.echo.Logger.Fatal(s.echo.Start(":1323"))
}

func (s *API) check() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}
}
