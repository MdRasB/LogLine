package auth

import (
	//"crypto"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	//"errors"
	"fmt"
)

const (
	APIKeyPrefix = "ll_live_"
)

func GenerateAPIKey() (string, string, error) {
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return "", "", fmt.Errorf("generate api key: %w", err)
	}

	apiKey := APIKeyPrefix + hex.EncodeToString(key)
	hkey := HashAPIKey(apiKey)

	return apiKey, hkey, nil
}

func HashAPIKey(s string) string {
	hash := sha256.Sum256([]byte(s))

	return hex.EncodeToString(hash[:])
}

func VerifyAPIKey(key, storedHash string) bool {
	computedHash := HashAPIKey(key)

	return subtle.ConstantTimeCompare(
		[]byte(computedHash),
		[]byte(storedHash),
	) == 1
}

