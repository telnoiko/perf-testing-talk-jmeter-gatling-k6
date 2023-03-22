package rest

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"task-service/app/store"
)

var (
	UserContextKey  = "user"
	secretKey       = "thisissecretkey"
	bearerExtractor = BearerExtractor{}
)

type Authorizer struct {
	store *store.Store
}

func NewAuth(store *store.Store) *Authorizer {
	return &Authorizer{store}
}

func (s *Authorizer) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := bearerExtractor.ExtractToken(c.Request())
		if err != nil {
			return c.String(http.StatusUnauthorized, "Could not extract Bearer token")
		}

		claims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			return err
		}

		id := claims.Claims.(*jwt.StandardClaims).Id

		// todo find user by id and token
		c.Set(UserContextKey, user)
		return next(c)
	}
}

func (s *Authorizer) GenerateJWT(id int) (string, error) {
	claims := jwt.StandardClaims{Id: strconv.Itoa(id), Issuer: "task-service"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(secretKey))
}
