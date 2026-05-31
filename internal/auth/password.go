package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCost = 12
	MinPasswordLength = 8
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password can't be empty")
	}

	if len(password) < 8 {
		return "", errors.New("password must be 8 character long")
	}

	hash , err := bcrypt.GenerateFromPassword(
		[]byte(password),
		BcryptCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword (password, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

