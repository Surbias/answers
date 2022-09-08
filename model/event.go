package model

type IEvent interface {
	Type() string
	Key() string
}

type AnswerEvent struct {
	EventType string `json:"event" bson:"event"`
	Timestamp int64  `json:"-" bson:"timestamp"`
	Data      Answer `json:"data" bson:"data"`
}

func (e AnswerEvent) Type() string {
	return e.EventType
}

func (e AnswerEvent) Key() string {
	return e.Data.Key
}
