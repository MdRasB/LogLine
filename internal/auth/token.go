package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

const (
	SessionTokenPrefix = "ll_sess_"
)

func GenerateSessionToken() (string, string, error) {
	token := make([]byte, 32)

	if _, err := rand.Read(token); err != nil {
		return "", "", fmt.Errorf("generate session token: %w", err)
	}

	sessionToken := SessionTokenPrefix + hex.EncodeToString(token)
	htoken := HashSessionToken(sessionToken)

	return sessionToken, htoken, nil
}

func HashSessionToken(s string) string {
	hash := sha256.Sum256([]byte(s))

	return hex.EncodeToString(hash[:])
}

func VerifySessionToken(token, storedHash string) bool {
	computedHash := HashSessionToken(token)

	return subtle.ConstantTimeCompare(
		[]byte(computedHash),
		[]byte(storedHash),
	) == 1
}