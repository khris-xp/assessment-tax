package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/khris-xp/assessment-tax/controllers"
)

func TaxRoutes(e *echo.Echo) {

	tCl := controllers.TaxController{}

	e.POST("/tax/calculations", tCl.CalculateTax)
	e.POST("tax/calculations/upload-csv", tCl.TaxCalculateFormCsv)
}
