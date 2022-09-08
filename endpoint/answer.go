package endpoint

import (
	"bequest/answers/common"
	"bequest/answers/model"
	"bequest/answers/repository"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func Answer(c echo.Context) error {
	answerKey := c.Param("key")

	answers := c.Get("answers").(repository.IAnswerRepository)
	answer, err := answers.GetAnswerByKey(answerKey)
	if err != nil {
		switch err.Error() {
		case fmt.Sprintf(common.ERR_ANSWER_NOT_FOUND, answerKey):
			c.Logger().Debug(fmt.Sprintf(common.ERR_ANSWER_NOT_FOUND, answerKey))
			return c.NoContent(http.StatusNotFound)
		default:
			c.Logger().Error(fmt.Sprintf(common.ERR_ACTION_ANSWER_KEY, "get", answerKey, err))
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, answer)
}

func CreateAnswer(c echo.Context) error {
	var answer model.Answer
	err := c.Bind(&answer)
	if err != nil {
		c.Logger().Debug(fmt.Sprintf(common.ERR_ACTION_ANSWER_INVALID_PAYLOAD, "create"))
		return c.NoContent(http.StatusBadRequest)
	}

	answers := c.Get("answers").(repository.IAnswerRepository)
	createdAnswer, err := answers.CreateAnswer(answer)
	if err != nil {
		switch err.Error() {
		case fmt.Sprintf(common.ERR_ANSWER_CONFLICT, answer.Key):
			c.Logger().Debug(fmt.Sprintf(common.ERR_ANSWER_CONFLICT, answer.Key))
			return c.NoContent(http.StatusConflict)
		default:
			c.Logger().Error(fmt.Sprintf(common.ERR_ACTION_ANSWER_KEY, "create", answer.Key, err))
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	events := c.Get("events").(repository.IEventRepository)
	eventType := "create"
	eventKey := createdAnswer.Key
	err = events.CreateEventByKey(model.AnswerEvent{
		Data:      createdAnswer,
		EventType: eventType,
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		c.Logger().Warn(fmt.Sprintf(common.ERR_EVENT_NOT_CREATED, eventType, eventKey, err))
	}
	return c.JSON(http.StatusCreated, createdAnswer)
}

func UpdateAnswer(c echo.Context) error {
	answerKey := c.Param("key")

	var update model.UpdateAnswer
	err := c.Bind(&update)
	if err != nil {
		c.Logger().Debug(fmt.Sprintf(common.ERR_ACTION_ANSWER_INVALID_PAYLOAD, "update"))
		return c.NoContent(http.StatusBadRequest)
	}

	answers := c.Get("answers").(repository.IAnswerRepository)
	updatedAnswer, err := answers.UpdateAnswerByKey(answerKey, update)
	if err != nil {
		switch err.Error() {
		case fmt.Sprintf(common.ERR_ANSWER_NOT_FOUND, answerKey):
			return c.NoContent(http.StatusNotFound)
		case common.ERR_ANSWER_VALUE_EMPTY:
			return c.JSON(http.StatusBadGateway, ErrorMessage{Message: common.ERR_ANSWER_VALUE_EMPTY})
		default:
			c.Logger().Error(fmt.Sprintf(common.ERR_ACTION_ANSWER_KEY, "update", answerKey, err))
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	events := c.Get("events").(repository.IEventRepository)
	eventType := "update"
	eventKey := answerKey
	err = events.CreateEventByKey(model.AnswerEvent{
		Data:      updatedAnswer,
		EventType: eventType,
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		c.Logger().Warn(fmt.Sprintf(common.ERR_EVENT_NOT_CREATED, eventType, eventKey, err))
	}

	return c.JSON(http.StatusOK, updatedAnswer)
}

func DeleteAnswer(c echo.Context) error {
	answerKey := c.Param("key")
	answers := c.Get("answers").(repository.IAnswerRepository)

	deletedAnswers, err := answers.DeleteAnswerByKey(answerKey)
	if err != nil {
		c.Logger().Error(fmt.Sprintf(common.ERR_ACTION_ANSWER_KEY, "delete", answerKey, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	events := c.Get("events").(repository.IEventRepository)
	eventType := "delete"
	eventKey := answerKey

	for _, deletedAnswer := range deletedAnswers {
		err = events.CreateEventByKey(model.AnswerEvent{
			Data:      deletedAnswer,
			EventType: eventType,
			Timestamp: time.Now().Unix(),
		})
		if err != nil {
			c.Logger().Warn(fmt.Sprintf(common.ERR_EVENT_NOT_CREATED, eventType, eventKey, err))
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func AnswerHistory(c echo.Context) error {
	answerKey := c.Param("key")

	events := c.Get("events").(repository.IEventRepository)
	foundEvents, err := events.GetEventsByKey(answerKey)
	if err != nil {
		switch err.Error() {
		case fmt.Sprintf(common.ERR_EVENT_NOT_FOUND, answerKey):
			c.Logger().Debug(fmt.Sprintf(common.ERR_EVENT_NOT_FOUND, answerKey))
			return c.NoContent(http.StatusNotFound)
		default:
			c.Logger().Error(fmt.Sprintf(common.ERR_EVENT_ANSWER_KEY, answerKey, err))
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return c.JSON(http.StatusOK, foundEvents)
}
