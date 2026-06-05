package auth

import (
	"errors"
	"time"

	"github.com/MdRasB/LogLine/internal/model"
)

const SessionDuration = 24 * time.Hour

var ErrSessionExpired = errors.New("session expired")

func (s *Service) ValidateSession(token string) (model.Session, error) {
	hash := HashSessionToken(token)

	session, err := s.sessions.GetSessionByTokenHash(hash)
	if err != nil {
		return model.Session{}, err
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.sessions.DeleteSession(session.ID)

		return model.Session{}, ErrSessionExpired
	}

	return session, nil
}
