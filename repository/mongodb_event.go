package repository

import (
	"bequest/answers/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	COLLECTION_EVENTS = "events"

	// ANSWER_KEY_PATH is the path to the answer key
	ANSWER_KEY_PATH = "data.key"
)

type MongoDBEvent struct {
	collection *mongo.Collection
}

func NewMongoDBEvent(db *mongo.Database) IEventRepository {
	return MongoDBEvent{collection: db.Collection(COLLECTION_EVENTS)}
}

func (e MongoDBEvent) CreateEventByKey(event model.IEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
	defer cancel()
	_, err := e.collection.InsertOne(ctx, event)
	return err
}

func (e MongoDBEvent) GetEventsByKey(key string) (events []model.IEvent, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
	defer cancel()

	cursor, err := e.collection.Find(ctx, bson.D{primitive.E{Key: ANSWER_KEY_PATH, Value: key}})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var event model.AnswerEvent
		err = cursor.Decode(&event)
		if err != nil {
			return
		}
		events = append(events, event)
	}
	return
}
