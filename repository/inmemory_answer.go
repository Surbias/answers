package repository

import (
	"bequest/answers/common"
	"bequest/answers/model"
	"errors"
	"fmt"
	"sort"
	"time"
)

type InMemoryAnswer struct {
	db map[string][]model.Answer
}

func NewInMemoryAnswer() IAnswerRepository {
	return InMemoryAnswer{db: make(map[string][]model.Answer)}
}

func (a InMemoryAnswer) recentAnswer(answers []model.Answer) model.Answer {
	if len(answers) == 0 {
		return model.Answer{}
	}

	sort.Slice(answers, func(p, q int) bool {
		return answers[p].Timestamp > answers[q].Timestamp
	})

	return answers[len(answers)-1]
}

func (a InMemoryAnswer) GetAnswerByKey(key string) (model.Answer, error) {
	answers, ok := a.db[key]
	if !ok {
		return model.Answer{}, fmt.Errorf(common.ERR_ANSWER_NOT_FOUND, key)
	}

	return a.recentAnswer(answers), nil
}

func (a InMemoryAnswer) CreateAnswer(answer model.Answer) (model.Answer, error) {
	answers, ok := a.db[answer.Key]
	if !ok {
		a.db[answer.Key] = []model.Answer{}
	}

	for _, currentAnswer := range answers {
		if currentAnswer.Value == answer.Value {
			return model.Answer{}, fmt.Errorf(common.ERR_ANSWER_CONFLICT, answer.Key)
		}
	}

	newAnswer := model.Answer{
		Key:       answer.Key,
		Timestamp: time.Now().Unix(),
		Value:     answer.Value,
	}
	a.db[newAnswer.Key] = append(a.db[newAnswer.Key], newAnswer)
	return newAnswer, nil
}

func (a InMemoryAnswer) DeleteAnswerByKey(key string) ([]model.Answer, error) {
	answers, ok := a.db[key]
	if ok {
		delete(a.db, key)
	}
	return answers, nil
}

func (a InMemoryAnswer) UpdateAnswerByKey(key string, updateAnswer model.UpdateAnswer) (model.Answer, error) {
	if updateAnswer.Value == "" {
		return model.Answer{}, errors.New(common.ERR_ANSWER_VALUE_EMPTY)
	}
	answers, ok := a.db[key]
	if !ok {
		return model.Answer{}, fmt.Errorf(common.ERR_ANSWER_NOT_FOUND, key)
	}

	for _, answer := range answers {
		if answer.Value == updateAnswer.Value {
			return answer, nil
		}
	}

	newAnswer := model.Answer{
		Key:       key,
		Timestamp: time.Now().Unix(),
		Value:     updateAnswer.Value,
	}

	a.db[key] = append(a.db[key], newAnswer)

	return newAnswer, nil
}
