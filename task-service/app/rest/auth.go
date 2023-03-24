package rest

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
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
		tokenString, err := bearerExtractor.ExtractToken(c.Request())
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

		user, err := s.store.FindByEmail(email)
		if err != nil {
			return c.String(http.StatusUnauthorized, "Could not find user with token")
		}

		c.Set(UserContextKey, user)
		return next(c)
	}
}

func (s *Authorizer) GenerateJWT(email string) (string, error) {
	bits := make([]byte, 12)
	_, err := rand.Read(bits)
	if err != nil {
		panic(err)
	}

	claims := jwt.MapClaims{
		"email": email,
		"iss":   "task-service",
		"jti":   base64.StdEncoding.EncodeToString(bits),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(secretKey))
}
