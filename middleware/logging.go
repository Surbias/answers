package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func SetLogLevel(level log.Lvl) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().SetLevel(level)
			return next(c)
		}
	}
}
