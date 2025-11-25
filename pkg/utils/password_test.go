package utils

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("hash password failed: %v", err)
	}

	if err := CheckPassword(hash, "secret"); err != nil {
		t.Fatalf("check password failed: %v", err)
	}

	if err := CheckPassword(hash, "wrong"); err == nil {
		t.Fatalf("expected error for wrong password")
	}
}
