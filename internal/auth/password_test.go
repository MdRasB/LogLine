package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "secret123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected hash to be generated")
	}

	if hash == password {
		t.Fatal("hash should not equal original password")
	}
}

func TestCheckPasswordSuccess(t *testing.T) {
	password := "secret123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	err = VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("expected password verification to succeed: %v", err)
	}
}

func TestCheckPasswordFail(t *testing.T) {
	password := "secret123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	err = VerifyPassword("wrong-password", hash)

	if err == nil {
		t.Fatal("expected verification to fail")
	}
}