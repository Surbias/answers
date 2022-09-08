package server

import (
	"bequest/answers/endpoint"

	"bequest/answers/repository"
	"fmt"

	m "bequest/answers/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type IServer interface {
	Listen(host string, port int) error
}

type server struct {
	e *echo.Echo
}

func (s server) Listen(host string, port int) error {
	return s.e.Start(fmt.Sprintf("%s:%d", host, port))
}

func NewServer(answers repository.IAnswerRepository, events repository.IEventRepository) IServer {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())

	answer := e.Group("/answers")
	answer.Use(m.SetLogLevel(log.DEBUG))
	answer.Use(m.Answers(answers))
	answer.Use(m.Events(events))

	answer.POST("", endpoint.CreateAnswer)
	answer.GET("/:key", endpoint.Answer)
	answer.PATCH("/:key", endpoint.UpdateAnswer)
	answer.DELETE("/:key", endpoint.DeleteAnswer)
	answer.GET("/:key/history", endpoint.AnswerHistory)

	return server{e: e}
}
