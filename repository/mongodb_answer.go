package repository

import (
	"bequest/answers/common"
	"bequest/answers/model"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	COLLECTION_ANSWERS = "answers"
	ANSWER_KEY         = "key"
)

type MongoDBAnswer struct {
	collection *mongo.Collection
}

func NewMongoDBAnswer(db *mongo.Database) IAnswerRepository {
	return MongoDBAnswer{collection: db.Collection(COLLECTION_ANSWERS)}
}

func (a MongoDBAnswer) filterByAnswerKey(value string) bson.D {
	return bson.D{primitive.E{Key: ANSWER_KEY, Value: value}}
}

func (a MongoDBAnswer) GetAnswerByKey(key string) (answer model.Answer, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
	defer cancel()

	filter := a.filterByAnswerKey(key)
	options := options.FindOne().SetSort(primitive.E{Key: "timestamp", Value: -1})
	result := a.collection.FindOne(ctx, filter, options).Decode(&answer)

	if result.Error() != "" {
		err = fmt.Errorf("mongodb error finding answer by key %q: %s", key, result.Error())
	}
	return
}

func (a MongoDBAnswer) CreateAnswer(answer model.Answer) (createdAnswer model.Answer, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
	defer cancel()

	filter := a.filterByAnswerKey(answer.Key)
	count, err := a.collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	if count != 0 {
		err = fmt.Errorf(common.ERR_ANSWER_CONFLICT, answer.Key)
		return
	}

	createdAnswer = model.Answer{
		Timestamp: time.Now().Unix(),
		Key:       answer.Key,
		Value:     answer.Value,
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
	defer cancel2()
	_, err = a.collection.InsertOne(ctx2, createdAnswer)
	if err != nil {
		return model.Answer{}, nil
	}
	return answer, nil
}

func (a MongoDBAnswer) DeleteAnswerByKey(key string) (answers []model.Answer, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
	defer cancel()

	filter := a.filterByAnswerKey(key)
	cursor, err := a.collection.Find(ctx, filter)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var answer model.Answer
		err = cursor.Decode(&answer)
		if err != nil {
			return
		}
		answers = append(answers, answer)
	}

	_, err = a.collection.DeleteMany(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (a MongoDBAnswer) UpdateAnswerByKey(key string, updateAnswer model.UpdateAnswer) (answer model.Answer, err error) {
	if updateAnswer.Value == "" {
		err = errors.New(common.ERR_ANSWER_VALUE_EMPTY)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
	defer cancel()

	filter := a.filterByAnswerKey(key)
	count, err := a.collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	if count == 0 {
		err = fmt.Errorf(common.ERR_ANSWER_NOT_FOUND, key)
		return
	}

	answer = model.Answer{
		Timestamp: time.Now().Unix(),
		Key:       key,
		Value:     updateAnswer.Value,
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
	defer cancel2()
	_, err = a.collection.InsertOne(ctx2, answer)
	if err != nil {
		answer = model.Answer{}
		return
	}
	return
}
