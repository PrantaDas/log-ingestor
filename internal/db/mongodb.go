package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDB struct {
	MongoClient    *mongo.Client
	CollectionName string
}

type Repository interface {
	InsertToDB(ctx context.Context) error
}

func ConnectToMongoDB(ctx context.Context, uri string, collName string) (*MongoDB, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return &MongoDB{
		MongoClient:    client,
		CollectionName: collName,
	}, nil
}
