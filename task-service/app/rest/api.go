package rest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"task-service/app/store"
)

type API struct {
	echo       *echo.Echo
	authorizer *Authorizer
	store      *store.Store
}

type SanitizedUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserApiResponse struct {
	User  SanitizedUser `json:"user"`
	Token string        `json:"token"`
}

type UserApiResponse2 struct {
	User  *store.User `json:"user"`
	Token string      `json:"token"`
}

func New(auth *Authorizer, store *store.Store) *API {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Recover())

	api := API{e, auth, store}

	return &api
}

func (s *API) Start() {
	s.echo.GET("/", s.check())
	s.echo.POST("/users", s.save())
	s.echo.POST("/login", s.login())
	//e.POST("/logout", rest.check(), s.authorizer.Authorize)

	//e.GET("/tasks", rest.check(), s.authorizer.Authorize)
	//e.GET("/tasks/:id", rest.check(), s.authorizer.Authorize)
	//e.POST("/tasks", rest.check(), s.authorizer.Authorize)
	//e.PUT("/tasks/:id", rest.check(), s.authorizer.Authorize)
	//e.DELETE("/tasks/:id", rest.check(), s.authorizer.Authorize)

	s.echo.Logger.Fatal(s.echo.Start(":1323"))
}

func (s *API) check() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
}

func (s *API) save() func(c echo.Context) error {
	return func(c echo.Context) error {
		user := store.User{}
		c.Bind(&user)

		id, err := s.store.Save(&user)
		if err != nil {
			fmt.Println(os.Stderr, err)
			return c.String(http.StatusBadRequest, "Unable to save user")
		}

		jwt, err := s.authorizer.GenerateJWT(id)
		if err != nil {
			fmt.Println(os.Stderr, err)
			return c.String(http.StatusInternalServerError, "Unable to generage JWT")
		}

		response := UserApiResponse{User: SanitizedUser{user.Name, user.Email}, Token: jwt}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
}

func (s *API) login() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(UserContextKey).(*store.User)
		if user == nil {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		foundUser, err := s.store.FindByCredentials(user.Email, user.Password)
		if err != nil {
			return err
		}

		jwt, err := s.authorizer.GenerateJWT(foundUser.ID)
		if err != nil {
			return err
		}

		response := UserApiResponse2{User: user, Token: jwt}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
}
