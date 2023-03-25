package rest

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"task-service/app/store"
)

type user struct {
	store      *store.Store
	authorizer *Authorizer
	logger     echo.Logger
}

type SanitizedUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserApiResponse struct {
	User  SanitizedUser `json:"user"`
	Token string        `json:"token"`
}

func (s *user) create() func(c echo.Context) error {
	return func(c echo.Context) error {
		user := store.User{}
		c.Bind(&user)
		c.Logger().Infof("create: received user: %v\n", user)

		jwt, err := s.authorizer.GenerateJWT(user.Email)
		user.Tokens = append(user.Tokens, jwt)
		if err != nil {
			s.logger.Infof("Unable to generage JWT: %v\n", err)
			return c.String(http.StatusInternalServerError, "Unable to generage JWT")
		}

		_, err = s.store.Save(&user)
		if err != nil {
			s.logger.Infof("Unable to create user: %v\n", err)
			return c.String(http.StatusBadRequest, "Unable to create user")
		}

		response := UserApiResponse{User: SanitizedUser{user.Name, user.Email}, Token: jwt}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
}

func (s *user) login() echo.HandlerFunc {
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
			s.logger.Infof("Unable to generage JWT: %v\n", err)
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

func (s *user) logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(UserContextKey).(*store.User)
		c.Logger().Infof("logout: received user: %v\n", user)

		err := s.store.DeleteToken(user.ID, user.Tokens[0])
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, "")
	}
}
