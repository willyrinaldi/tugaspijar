package main

import (
	"tugaspijar/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))

	h := handler.NewHandler()
	e.GET("/", h.Index)
	e.POST("/", h.Create)
	e.PUT("/:id", h.Update)
	e.GET("/:id", h.SimpleOne)
	e.Logger.Fatal(e.Start(":1323"))
}