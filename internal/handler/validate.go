package handler

import (
	"errors"
	"time"

	"github.com/MdRasB/LogLine/internal/model"
)

var ValidateLevel = map[string]bool{
	"error": true, "warn": true,
	"info":  true, "debug": true,
	"fatal": true,
}

func validateLevels(level string) error {
	if level == "" {
		return errors.New("level is required")
	} else if !ValidateLevel[level] {
		return errors.New("invalid log level")
	}

	return nil
}

func Validate(log model.Logs) error {
	if err := validateLevels(log.Level); err != nil {
		return err
	}

	if log.Message == "" {
		return errors.New("message is required")
	}

	if log.Service == "" {
		return errors.New("service is required")
	}

	if log.Timestamp == "" {
		return errors.New("timestamp is required")
	}

	if _, err := time.Parse(time.RFC3339, log.Timestamp); err != nil {
		return errors.New("invalid timestamp format")
	}


	return nil
}