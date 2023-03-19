package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"tasks-service/app/middleware"
)

func main() {

	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.AuthorizerMiddleware)

	e.GET("/", check())

	e.POST("/users", check())
	e.POST("/login", check())
	e.POST("/logout", check())

	e.GET("/tasks", check())
	e.GET("/tasks/:id", check())
	e.POST("/tasks", check())
	e.PUT("/tasks/:id", check())
	e.DELETE("/tasks/:id", check())

	e.Logger.Fatal(e.Start(":1323"))
}
