package auth

import (
	"testing"
	"time"

	"go-mdbook/internal/config"
)

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if !CheckPassword(hash, "secret") {
		t.Fatalf("expected password to match")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatalf("expected password to fail")
	}
}

func TestGenerateAndParseToken(t *testing.T) {
	cfg := config.Config{JWTSecret: "test", TokenTTL: time.Hour}
	token, err := GenerateToken(cfg, "user123", "admin")
	if err != nil {
		t.Fatalf("token error: %v", err)
	}
	claims, err := ParseToken(cfg, token)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if claims.UserID != "user123" || claims.Role != "admin" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}
