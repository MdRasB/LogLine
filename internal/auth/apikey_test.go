package auth

import "testing"

func TestGenerateAPIKey(t *testing.T) {
	key, _, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if key == "" {
		t.Fatal("expected generated api key")
	}
}

func TestGenerateAPIKeyUnique(t *testing.T) {
	key1, _,  err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("generate key1: %v", err)
	}

	key2, _, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("generate key2: %v", err)
	}

	if key1 == key2 {
		t.Fatal("generated api keys should be unique")
	}
}

func TestHashAPIKey(t *testing.T) {
	key := "ll_live_test_key"

	hash := HashAPIKey(key)

	if hash == "" {
		t.Fatal("expected hash")
	}

	if hash == key {
		t.Fatal("hash should not equal original key")
	}

	if len(hash) != 64 {
		t.Fatalf("expected hash length 64, got %d", len(hash))
	}
}

func TestHashAPIKeyDeterministic(t *testing.T) {
	key := "ll_live_test_key"

	hash1 := HashAPIKey(key)
	hash2 := HashAPIKey(key)

	if hash1 != hash2 {
		t.Fatal("same key should always produce same hash")
	}
}

func TestVerifyAPIKeySuccess(t *testing.T) {
	key := "ll_live_test_key"

	hash := HashAPIKey(key)

	ok := VerifyAPIKey(key, hash)

	if !ok {
		t.Fatal("expected api key verification to succeed")
	}
}

func TestVerifyAPIKeyFail(t *testing.T) {
	key := "ll_live_test_key"

	hash := HashAPIKey(key)

	ok := VerifyAPIKey("wrong-key", hash)

	if ok {
		t.Fatal("expected verification to fail")
	}
}