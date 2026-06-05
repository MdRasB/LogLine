package db

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidToken = errors.New("invalid token")
)