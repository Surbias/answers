package repository

import "bequest/answers/model"

//go:generate genmock
type IEventRepository interface {
	CreateEventByKey(event model.IEvent) error
	GetEventsByKey(key string) ([]model.IEvent, error)
}
