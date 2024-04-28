package routes

import (
	"github.com/khris-xp/assessment-tax/controllers"
	"github.com/khris-xp/assessment-tax/middlewares"
	"github.com/labstack/echo/v4"
)

func KReceiptRoutes(e *echo.Echo) {

	rCl := controllers.KReceiptController{}

	e.POST("/admin/deductions/k-receipt", rCl.CreateKReceipt, middlewares.AuthMiddleware)
}
