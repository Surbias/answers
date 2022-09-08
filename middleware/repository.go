package middleware

import (
	"bequest/answers/repository"

	"github.com/labstack/echo/v4"
)

func Answers(answers repository.IAnswerRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("answers", answers)
			return next(c)
		}
	}
}

func Events(events repository.IEventRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("events", events)
			return next(c)
		}
	}
}
