package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func check() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
}

func save() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
}
