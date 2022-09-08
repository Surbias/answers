package main

import (
	"bequest/answers/repository"
	"bequest/answers/server"
	"context"
	"flag"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var host string
	var port int

	flag.StringVar(&host, "host", "localhost", "the server host")
	flag.IntVar(&port, "port", 8081, "the server port")
	flag.Parse()

	var config struct {
		MongoDbUri  string `envconfig:"MONGODB_URI"`
		DbName      string `envconfig:"DB_NAME" default:"bequest"`
		UseInMemory bool   `envconfig:"USE_IN_MEMORY" default:"false"`
	}

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("could not extract environment variables: %v", err)
	}

	var answers repository.IAnswerRepository
	var events repository.IEventRepository

	if !config.UseInMemory {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDbUri))
		if err != nil {
			log.Fatalf("could not create MongoDB client: %v", err)
		}
		err = client.Connect(ctx)
		if err != nil {
			log.Fatalf("could not connect to MongoDB database: %v", err)
		}
		defer client.Disconnect(ctx)

		database := client.Database(config.DbName)
		answers = repository.NewMongoDBAnswer(database)
		events = repository.NewMongoDBEvent(database)
	} else {
		answers = repository.NewInMemoryAnswer()
		events = repository.NewInMemoryEvent()
	}

	s := server.NewServer(answers, events)
	if err := s.Listen(host, port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
