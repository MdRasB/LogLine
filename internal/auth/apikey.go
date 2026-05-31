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
	TestAppend = "ll_live_"
)

func GenerateAPIKey()(string, error) {
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("generate api key: %w\n", err)
	}

	return	APIKeyPrefix + hex.EncodeToString(key), nil
}

func HashAPIKey(s string) (string){
	hash := sha256.Sum256([]byte(s))

	return hex.EncodeToString(hash[:])
}

func VerifyAPIKey(key, storedHash string) bool {
	ComputedHash := HashAPIKey(key)

	return subtle.ConstantTimeCompare(
		[]byte(ComputedHash),
		[]byte(storedHash),
	) == 1
}