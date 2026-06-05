package db

import (
	"context"
	//"errors"
	"fmt"

	"github.com/MdRasB/LogLine/internal/model"
	"github.com/google/uuid"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (u *UserStore) CreateUser(user model.User) error {
	query := `
			Insert Into users (id, email, password_hash, created_at)
			Values ($1, $2, $3, $4)
	`

	_, err := u.db.Exec(
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

func (u *UserStore) GetUserByEmail(email string) (model.User, error) {
	query := `
			SELECT id, email, password_hash, created_at
			FROM users
			WHERE email = $1
	`
	var user model.User

	err := u.db.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt)

	if err != nil {
		return user, fmt.Errorf("getting user by email: %w", err)
	}

	return user, nil
}

func (u *UserStore) GetUserByID(id uuid.UUID) (model.User, error) {
	query := `
			SELECT id, email, password_hash, created_at
			FROM users
			WHERE id = $1
	`

	var user model.User

	err := u.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt)

	if err != nil {
		return user, fmt.Errorf("getting user by id: %w", err)
	}

	return user, nil
}