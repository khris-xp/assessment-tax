package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/khris-xp/assessment-tax/controllers"
	"github.com/khris-xp/assessment-tax/middlewares"
)

func TaxRoutes(e *echo.Echo) {

	tCl := controllers.TaxController{}

	e.POST("/tax/calculations", tCl.CalculateTax)
	e.POST("/admin/deductions/personal", tCl.CreatePersonalDeduction, middlewares.AuthMiddleware)
}
