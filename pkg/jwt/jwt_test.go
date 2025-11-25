package jwtutil

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken(1, "tester", "admin", "secret", 1)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	claims, err := ParseToken(token, "secret")
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}

	if claims.UserID != 1 || claims.Username != "tester" || claims.Role != "admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}

	if time.Until(claims.ExpiresAt.Time) <= 0 {
		t.Fatalf("token already expired")
	}
}

func TestValidateToken(t *testing.T) {
	token, err := GenerateToken(2, "user", "user", "secret", 1)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	if !ValidateToken(token, "secret") {
		t.Fatalf("expected valid token")
	}

	if ValidateToken(token, "wrong") {
		t.Fatalf("expected invalid token for wrong secret")
	}
}
