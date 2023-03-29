package rest

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"task-service/app/store"
)

type user struct {
	store  *store.Store
	token  *Token
	logger echo.Logger
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
		log.Printf("create: received user: %v\n", user)

		jwt, err := s.token.GenerateJWT(user.Email)
		user.Tokens = append(user.Tokens, jwt)
		if err != nil {
			log.Printf("Unable to generage JWT: %v\n", err)
			return c.String(http.StatusInternalServerError, "Unable to generage JWT")
		}

		err = s.store.User.Create(&user)
		if err != nil {
			log.Printf("Unable to create user: %v\n", err)
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
		log.Printf("login: received user: %v\n", user)

		foundUser, err := s.store.User.FindByEmail(user.Email)
		if err != nil {
			log.Printf("FindByEmail failed: %v\n", err)
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		log.Printf("foundUser: %v\n", foundUser)

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
		if err != nil {
			log.Printf("CompareHashAndPassword failed: %v\n", err)
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		jwt, err := s.token.GenerateJWT(foundUser.Email)
		if err != nil {
			log.Printf("Unable to generage JWT: %v\n", err)
			return c.String(http.StatusInternalServerError, "Unable to generage JWT")
		}
		err = s.store.User.UpdateToken(foundUser.ID, jwt)
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
		log.Printf("logout: received user: %v\n", user)

		err := s.store.User.DeleteToken(user.ID, user.Tokens[0])
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, "")
	}
}

func (s *user) logoutAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(UserContextKey).(*store.User)
		log.Printf("logoutAll: received user: %v\n", user)

		err := s.store.User.DeleteAllTokens(user.ID)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, "")
	}
}
