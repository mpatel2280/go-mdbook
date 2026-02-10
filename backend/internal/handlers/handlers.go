package handlers

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go-mdbook/internal/auth"
	"go-mdbook/internal/config"
	"go-mdbook/internal/models"
	"go-mdbook/internal/services"
	"go-mdbook/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	cfg    config.Config
	client *mongo.Client
}

func New(cfg config.Config, client *mongo.Client) *Handler {
	return &Handler{cfg: cfg, client: client}
}

func (h *Handler) users() *mongo.Collection {
	return h.client.Database(h.cfg.MongoDB).Collection("users")
}

func (h *Handler) books() *mongo.Collection {
	return h.client.Database(h.cfg.MongoDB).Collection("books")
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var user models.User
	if err := h.users().FindOne(h.cfg.Context(), bson.M{"email": strings.ToLower(req.Email), "active": true}).Decode(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := auth.GenerateToken(h.cfg, user.ID.Hex(), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "role": user.Role, "email": user.Email})
}

func (h *Handler) Me(c *gin.Context) {
	userID := c.GetString("userId")
	objID, _ := primitive.ObjectIDFromHex(userID)
	var user models.User
	if err := h.users().FindOne(h.cfg.Context(), bson.M{"_id": objID}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) ListUsers(c *gin.Context) {
	cur, err := h.users().Find(h.cfg.Context(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
		return
	}
	defer cur.Close(h.cfg.Context())

	users := []models.User{}
	if err := cur.All(h.cfg.Context(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode"})
		return
	}
	c.JSON(http.StatusOK, users)
}

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if req.Role == "" {
		req.Role = "reader"
	}
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password required"})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash"})
		return
	}
	user := models.User{Email: strings.ToLower(req.Email), PasswordHash: hash, Role: req.Role, Active: true}
	_, err = h.users().InsertOne(h.cfg.Context(), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

type updateUserRequest struct {
	Role   *string `json:"role"`
	Active *bool   `json:"active"`
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	update := bson.M{}
	if req.Role != nil {
		update["role"] = *req.Role
	}
	if req.Active != nil {
		update["active"] = *req.Active
	}
	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no changes"})
		return
	}
	_, err = h.users().UpdateByID(h.cfg.Context(), objID, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	_, err = h.users().DeleteOne(h.cfg.Context(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) ListBooks(c *gin.Context) {
	cur, err := h.books().Find(h.cfg.Context(), bson.M{"active": true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
		return
	}
	defer cur.Close(h.cfg.Context())
	books := []models.Book{}
	if err := cur.All(h.cfg.Context(), &books); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *Handler) GetBook(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var book models.Book
	if err := h.books().FindOne(h.cfg.Context(), bson.M{"_id": objID, "active": true}).Decode(&book); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

type createBookRequest struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

func (h *Handler) CreateBook(c *gin.Context) {
	var req createBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title required"})
		return
	}
	slug := req.Slug
	if slug == "" {
		slug = utils.Slugify(req.Title)
	}

	sourceDir := filepath.Join(h.cfg.BooksRoot, slug)
	buildDir := filepath.Join(h.cfg.BooksBuildRoot, slug)

	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create source dir"})
		return
	}
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create build dir"})
		return
	}

	book := models.Book{Title: req.Title, Slug: slug, SourceDir: sourceDir, BuildDir: buildDir, Active: true}
	res, err := h.books().InsertOne(h.cfg.Context(), book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug already exists"})
		return
	}
	book.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, book)
}

type updateBookRequest struct {
	Title  *string `json:"title"`
	Active *bool   `json:"active"`
}

func (h *Handler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	update := bson.M{}
	if req.Title != nil {
		update["title"] = *req.Title
	}
	if req.Active != nil {
		update["active"] = *req.Active
	}
	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no changes"})
		return
	}
	_, err = h.books().UpdateByID(h.cfg.Context(), objID, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	_, err = h.books().DeleteOne(h.cfg.Context(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) BuildBook(c *gin.Context) {
	book, ok := h.bookByID(c)
	if !ok {
		return
	}

	if err := services.BuildBook(book.SourceDir, book.BuildDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "built"})
}

func (h *Handler) UploadBook(c *gin.Context) {
	book, ok := h.bookByID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".zip") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .zip files are supported"})
		return
	}

	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, file.Filename)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save upload"})
		return
	}
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if book.SourceDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "missing source directory"})
		return
	}

	if err := os.RemoveAll(book.SourceDir); err != nil && !errors.Is(err, os.ErrNotExist) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear source directory"})
		return
	}
	if err := os.MkdirAll(book.SourceDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to recreate source directory"})
		return
	}
	if err := os.RemoveAll(book.BuildDir); err != nil && !errors.Is(err, os.ErrNotExist) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear build directory"})
		return
	}

	if err := services.ExtractZip(tmpPath, book.SourceDir); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "uploaded"})
}

func (h *Handler) BookContent(c *gin.Context) {
	book, ok := h.bookByID(c)
	if !ok {
		return
	}

	filepathParam := strings.TrimPrefix(c.Param("filepath"), "/")
	if filepathParam == "" {
		filepathParam = "index.html"
	}
	cleanPath := filepath.Clean(filepathParam)
	full := filepath.Join(book.BuildDir, cleanPath)
	buildRoot := filepath.Clean(book.BuildDir) + string(os.PathSeparator)
	if !strings.HasPrefix(filepath.Clean(full), buildRoot) && filepath.Clean(full) != filepath.Clean(book.BuildDir) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}
	c.File(full)
}

func (h *Handler) bookByID(c *gin.Context) (models.Book, bool) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return models.Book{}, false
	}
	var book models.Book
	if err := h.books().FindOne(h.cfg.Context(), bson.M{"_id": objID}).Decode(&book); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return models.Book{}, false
	}
	return book, true
}
