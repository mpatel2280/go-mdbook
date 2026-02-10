package db

import (
	"context"

	"go-mdbook/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(cfg config.Config) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}
	return client, nil
}

func collection(cfg config.Config, client *mongo.Client, name string) *mongo.Collection {
	return client.Database(cfg.MongoDB).Collection(name)
}
