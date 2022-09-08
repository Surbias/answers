package repository

import "bequest/answers/model"

//go:generate genmock
type IAnswerRepository interface {
	GetAnswerByKey(key string) (model.Answer, error)
	CreateAnswer(answer model.Answer) (model.Answer, error)
	DeleteAnswerByKey(key string) ([]model.Answer, error)
	UpdateAnswerByKey(key string, updateAnswer model.UpdateAnswer) (model.Answer, error)
}
