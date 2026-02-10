package main

import (
	"log"

	"go-mdbook/internal/config"
	"go-mdbook/internal/db"
	"go-mdbook/internal/handlers"
	"go-mdbook/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	client, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer func() {
		_ = client.Disconnect(cfg.Context())
	}()

	if err := db.EnsureIndexes(cfg, client); err != nil {
		log.Fatalf("ensure indexes: %v", err)
	}
	if err := db.EnsureAdmin(cfg, client); err != nil {
		log.Fatalf("ensure admin: %v", err)
	}

	r := gin.Default()
	r.Use(middleware.CORS())

	h := handlers.New(cfg, client)

	api := r.Group("/api")
	{
		api.POST("/auth/register", h.Register)
		api.POST("/auth/login", h.Login)
	}

	protected := api.Group("")
	protected.Use(middleware.Auth(cfg))
	{
		protected.GET("/me", h.Me)
		protected.GET("/books", h.ListBooks)
		protected.GET("/books/:id", h.GetBook)
		protected.GET("/books/:id/content/*filepath", h.BookContent)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.Auth(cfg), middleware.RequireRole("admin"))
	{
		admin.GET("/users", h.ListUsers)
		admin.POST("/users", h.CreateUser)
		admin.PATCH("/users/:id", h.UpdateUser)
		admin.DELETE("/users/:id", h.DeleteUser)

		admin.POST("/books", h.CreateBook)
		admin.PATCH("/books/:id", h.UpdateBook)
		admin.DELETE("/books/:id", h.DeleteBook)
		admin.POST("/books/:id/upload", h.UploadBook)
		admin.POST("/books/:id/build", h.BuildBook)
	}

	log.Printf("listening on %s", cfg.APIAddr)
	if err := r.Run(cfg.APIAddr); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
