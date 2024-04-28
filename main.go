package main

import (
	"net/http"

	"github.com/khris-xp/assessment-tax/config"
	"github.com/khris-xp/assessment-tax/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	config.DatabaseInit()
	gorm := config.DB()

	dbGorm, err := gorm.DB()
	if err != nil {
		panic(err)
	}

	routes.TaxRoutes(e)
	port := config.EnvPort()
	dbGorm.Ping()

	e.Logger.Fatal(e.Start(":" + port))
}
