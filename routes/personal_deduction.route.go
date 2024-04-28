package routes

import (
	"github.com/khris-xp/assessment-tax/controllers"
	"github.com/khris-xp/assessment-tax/middlewares"
	"github.com/labstack/echo/v4"
)

func PersonalDeductionRoutes(e *echo.Echo) {

	pCl := controllers.PersonalDeductionController{}

	e.POST("/admin/deductions/personal", pCl.CreatePersonalDeduction, middlewares.AuthMiddleware)
}
