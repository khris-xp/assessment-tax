package middlewares

import (
	"github.com/khris-xp/assessment-tax/config"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if config.EnvAdminUsername() == "adminTax" && config.EnvAdminPassword() == "admin!" {
			return next(c)
		}

		return c.JSON(401, "Unauthorized")
	}
}
