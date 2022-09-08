package endpoint_test

import (
	"bequest/answers/endpoint"
	"bequest/answers/model"
	"bequest/answers/repository/mock"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAnswerStatusOK(t *testing.T) {
	keyParam := "somekey"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/answers/%s", keyParam), nil)
	rec := httptest.NewRecorder()

	tt, _ := time.Parse(time.RFC3339, "2022-09-10T13:30:00Z")
	expectedAnswer := model.Answer{
		Timestamp: tt.Unix(),
		Key:       keyParam,
		Value:     "some-value",
	}

	answers := mock.NewIAnswerRepository_Mock(t)
	answers.ReturnMockErrorAsResult = true
	answers.On_GetAnswerByKey().Rets(expectedAnswer, nil).Times(1)

	context := echo.New().NewContext(req, rec)
	context.Set("answers", answers)

	err := endpoint.Answer(context)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	body, err := ioutil.ReadAll(rec.Body)
	assert.NoError(t, err)

	var receivedAnswer model.Answer
	err = json.Unmarshal(body, &receivedAnswer)
	assert.NoError(t, err)

	assert.Equal(t, expectedAnswer.Key, receivedAnswer.Key)
	assert.Equal(t, expectedAnswer.Value, receivedAnswer.Value)

	// server does not retrieve timestamp in json payload
	assert.NotEqual(t, expectedAnswer.Timestamp, receivedAnswer.Timestamp)
}

func TestAnswerStatusInternalServerError(t *testing.T) {
	keyParam := "somekey"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/answers/%s", keyParam), nil)
	rec := httptest.NewRecorder()

	answers := mock.NewIAnswerRepository_Mock(t)
	answers.ReturnMockErrorAsResult = false
	answers.On_GetAnswerByKey().Rets(model.Answer{}, errors.New("something went wrong!")).Times(1)

	context := echo.New().NewContext(req, rec)
	context.Set("answers", answers)

	err := endpoint.Answer(context)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, &bytes.Buffer{}, rec.Body)
}
