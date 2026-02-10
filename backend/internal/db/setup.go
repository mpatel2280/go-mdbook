package db

import (
	"go-mdbook/internal/auth"
	"go-mdbook/internal/config"
	"go-mdbook/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndexes(cfg config.Config, client *mongo.Client) error {
	users := collection(cfg, client, "users")
	_, err := users.Indexes().CreateOne(cfg.Context(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	books := collection(cfg, client, "books")
	_, err = books.Indexes().CreateOne(cfg.Context(), mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func EnsureAdmin(cfg config.Config, client *mongo.Client) error {
	users := collection(cfg, client, "users")
	count, err := users.CountDocuments(cfg.Context(), bson.M{"email": cfg.AdminEmail})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := auth.HashPassword(cfg.AdminPassword)
	if err != nil {
		return err
	}

	_, err = users.InsertOne(cfg.Context(), models.User{
		Email:        cfg.AdminEmail,
		PasswordHash: hash,
		Role:         "admin",
		Active:       true,
	})
	return err
}
