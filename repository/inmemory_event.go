package repository

import (
	"bequest/answers/common"
	"bequest/answers/model"
	"fmt"
)

type InMemoryEvent struct {
	db map[string][]model.IEvent
}

func NewInMemoryEvent() IEventRepository {
	return InMemoryEvent{db: map[string][]model.IEvent{}}
}

func (e InMemoryEvent) CreateEventByKey(event model.IEvent) error {
	e.db[event.Key()] = append(e.db[event.Key()], event)
	return nil
}
func (e InMemoryEvent) GetEventsByKey(key string) ([]model.IEvent, error) {
	events, ok := e.db[key]
	if !ok {
		return nil, fmt.Errorf(common.ERR_EVENT_NOT_FOUND, key)
	}
	return events, nil
}
