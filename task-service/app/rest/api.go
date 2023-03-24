package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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

func New(auth *Authorizer, store *store.Store) *API {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Recover())

	api := API{e, auth, store}

	return &api
}

func (s *API) Start() {
	s.echo.GET("/", s.check())
	s.echo.POST("/users", s.create())
	s.echo.POST("/users/login", s.login())
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
		return c.String(http.StatusOK, "")
	}
}

func (s *API) create() func(c echo.Context) error {
	return func(c echo.Context) error {
		user := store.User{}
		c.Bind(&user)
		c.Logger().Infof("create: received user: %v\n", user)

		jwt, err := s.authorizer.GenerateJWT(user.Email)
		user.Tokens = append(user.Tokens, jwt)
		if err != nil {
			s.echo.Logger.Infof("Unable to generage JWT: %v\n", err)
			return c.String(http.StatusInternalServerError, "Unable to generage JWT")
		}

		_, err = s.store.Save(&user)
		if err != nil {
			s.echo.Logger.Infof("Unable to create user: %v\n", err)
			return c.String(http.StatusBadRequest, "Unable to create user")
		}

		response := UserApiResponse{User: SanitizedUser{user.Name, user.Email}, Token: jwt}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
}

func (s *API) login() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := store.User{}
		c.Bind(&user)
		c.Logger().Infof("login: received user: %v\n", user)

		foundUser, err := s.store.FindByEmail(user.Email)
		if err != nil {
			c.Logger().Warnf("FindByEmail failed: %v\n", err)
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		c.Echo().Logger.Infof("foundUser: %v\n", foundUser)

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
		if err != nil {
			c.Logger().Warnf("CompareHashAndPassword failed: %v\n", err)
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		jwt, err := s.authorizer.GenerateJWT(foundUser.Email)
		if err != nil {
			s.echo.Logger.Infof("Unable to generage JWT: %v\n", err)
			return c.String(http.StatusInternalServerError, "Unable to generage JWT")
		}
		err = s.store.UpdateToken(foundUser.ID, jwt)
		if err != nil {
			return err
		}

		response := UserApiResponse{User: SanitizedUser{foundUser.Name, foundUser.Email}, Token: jwt}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
}
