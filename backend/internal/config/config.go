package config

import (
	"context"
	"os"
	"time"
)

type Config struct {
	APIAddr        string
	MongoURI       string
	MongoDB        string
	JWTSecret      string
	TokenTTL       time.Duration
	AdminEmail     string
	AdminPassword  string
	BooksRoot      string
	BooksBuildRoot string
}

func Load() Config {
	return Config{
		APIAddr:        getEnv("API_ADDR", ":8080"),
		MongoURI:       getEnv("MONGO_URI", "mongodb://mongo:27017"),
		MongoDB:        getEnv("MONGO_DB", "mdbook"),
		JWTSecret:      getEnv("JWT_SECRET", "dev_secret_change_me"),
		TokenTTL:       getDuration("TOKEN_TTL", 24*time.Hour),
		AdminEmail:     getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "admin123"),
		BooksRoot:      getEnv("BOOKS_ROOT", "/data/books"),
		BooksBuildRoot: getEnv("BOOKS_BUILD_ROOT", "/data/build"),
	}
}

func (c Config) Context() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
