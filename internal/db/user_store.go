package db

import (
	"context"

	"github.com/MdRasB/LogLine/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db * pgxpool.Pool) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) UserCreate(user model.User) error {
	query := `
			Insert Into users (id, email, password_hash, created_at)
			Values ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(
		context.Background(),
		query,
		user.Id,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}