package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/MdRasB/LogLine/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionStore struct {
	db *pgxpool.Pool
}

func NewSessionStore(db *pgxpool.Pool) *SessionStore {
	return &SessionStore{
		db: db,
	}
}

func (s *SessionStore) CreateSession(session model.Session) error {
	query := `
			Insert Into sessions (id, user_id, token_hash, expires_at, created_at)
			Values ($1, $2, $3, $4, $5)
	`

	_, err := s.db.Exec(
		context.Background(),
		query,
		session.ID,
		session.UserID,
		session.TokenHash,
		session.ExpiresAt,
		session.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil 
}

func (s *SessionStore) GetSessionByTokenHash(tokenHash string) (model.Session, error) {
	query := `
			SELECT id, user_id, token_hash, expires_at, created_at
			From sessions
			Where token_hash = $1
	`

	var session model.Session

	err := s.db.QueryRow(
		context.Background(),
		query,
		tokenHash,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return session, ErrSessionNotFound
	}

	if err != nil {
		return session, fmt.Errorf("getting session by token: %w", err)

	}

	return session, nil
}

func (s *SessionStore) DeleteSession(id uuid.UUID) error {
	query := `
			DELETE FROM sessions
			WHERE id = $1
	`

	result, err := s.db.Exec(
		context.Background(),
		query,
		id,
	)

	if err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrSessionNotFound
	}

	return nil
}

func (s *SessionStore) DeleteExpiredSessions() error {
	query := `
			DELETE FROM sessions
			WHERE expires_at < Now()
	`

	_, err := s.db.Exec(
		context.Background(),
		query,
	)

	if err != nil {
		return fmt.Errorf("deleting expired session: %w", err)
	}

	return nil 
}