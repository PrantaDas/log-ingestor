package db

import (
	"context"
	"log"
	"log-ingester/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDB struct {
	MongoClient    *mongo.Client
	DbName         string
	CollectionName string
}

type Repository interface {
	InsertLogToDB(ctx context.Context) error
}

func ConnectToMongoDB(ctx context.Context, config *config.Config) (*MongoDB, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(config.MongoDBURI))
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

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
		CollectionName: config.CollectionName,
		DbName:         config.DBName,
	}, nil
}

func (m *MongoDB) InsertLogToDB(ctx context.Context, data string) error {
	collection := m.MongoClient.Database(m.DbName).Collection(m.CollectionName)

	document := bson.M{
		"foo": data,
	}
	if _, err := collection.InsertOne(ctx, document); err != nil {
		return err
	}
	return nil
}
