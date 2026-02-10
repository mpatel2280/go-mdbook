package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Role         string             `bson:"role" json:"role"`
	Active       bool               `bson:"active" json:"active"`
}

type Book struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Slug      string             `bson:"slug" json:"slug"`
	SourceDir string             `bson:"source_dir" json:"sourceDir"`
	BuildDir  string             `bson:"build_dir" json:"buildDir"`
	Active    bool               `bson:"active" json:"active"`
}
