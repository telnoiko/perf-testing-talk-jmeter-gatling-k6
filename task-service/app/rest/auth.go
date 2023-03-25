package rest

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"task-service/app/store"
)

var (
	UserContextKey = "user"
	secretKey      = "thisissecretkey"
)

type Authorizer struct {
	store *store.Store
	token *Token
}

func NewAuth(store *store.Store, token *Token) *Authorizer {
	return &Authorizer{store, token}
}

func (s *Authorizer) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString, err := s.token.ExtractToken(c.Request())
		if err != nil {
			return c.String(http.StatusUnauthorized, "Could not extract Bearer token")
		}

		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			return c.String(http.StatusUnauthorized, "Could not parse Bearer token")
		}

		var email string
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email = claims["email"].(string)
		}
		if err != nil {
			return c.String(http.StatusInternalServerError, "Could not parse id")
		}

		user, err := s.store.User.FindByToken(email, tokenString)
		if err != nil {
			return c.String(http.StatusUnauthorized, "Could not find user with token")
		}
		user.Tokens = []string{tokenString}

		c.Set(UserContextKey, user)
		return next(c)
	}
}
