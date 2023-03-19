package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func isAuthorized(c echo.Context) bool {
	// Your code to check if the user is authorized goes here
	return true
}

func AuthorizerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !isAuthorized(c) {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}
