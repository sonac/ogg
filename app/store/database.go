package store

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var ctx = context.TODO()

type Mongo struct {
	Client *mongo.Client
	UserCollection *mongo.Collection
	WordCollection *mongo.Collection
}

func NewDatabase() *Mongo {
	mg := Mongo{}
	mg.Connect()
	return &mg
}

func (m *Mongo) Connect() {
	mongoURI := os.Getenv("MONGO_URL")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	m.Client = client
	m.UserCollection = client.Database("ogg").Collection("users")
	m.WordCollection = client.Database("ogg").Collection("words")
}
